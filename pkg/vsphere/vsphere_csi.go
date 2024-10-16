package vsphere

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/gcfg.v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/legacy-cloud-providers/vsphere"
	"sigs.k8s.io/vsphere-csi-driver/v3/pkg/csi/service"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	cnstypes "github.com/vmware/govmomi/cns/types"
	vim25types "github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vslm"

	configclient "github.com/openshift/client-go/config/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cnsvolume "sigs.k8s.io/vsphere-csi-driver/v3/pkg/common/cns-lib/volume"
	cnsvsphere "sigs.k8s.io/vsphere-csi-driver/v3/pkg/common/cns-lib/vsphere"
	cnsconfig "sigs.k8s.io/vsphere-csi-driver/v3/pkg/common/config"
	"sigs.k8s.io/vsphere-csi-driver/v3/pkg/csi/service/common"
	"sigs.k8s.io/vsphere-csi-driver/v3/pkg/csi/service/logger"
)

const (
	clientOperatorName   = "cns-register"
	cloudCredSecretName  = "vmware-vsphere-cloud-credentials"
	csiDriverNamespace   = "openshift-cluster-csi-drivers"
	cloudConfigNamespace = "openshift-config"
	vsphereCsiConfig     = "vsphere-csi-config"
	csiConfigLocation    = "/tmp/vsphere-csi.conf"
	EnvVSphereCSIConfig  = "VSPHERE_CSI_CONFIG"
)

type VsphereCSIDebugger struct {
	config                   *vsphere.VSphereConfig
	kubeConfig               *rest.Config
	clientSet                kubernetes.Interface
	openshiftConfigClientSet configclient.Interface
	vCenter                  *cnsvsphere.VirtualCenter
}

func NewVSphereCSIDebugger(config *rest.Config) (*VsphereCSIDebugger, error) {
	service.Version = "1.0.1"
	t := &VsphereCSIDebugger{
		kubeConfig: config,
	}
	return t, nil
}

func (csi *VsphereCSIDebugger) createClients() {
	// create the clientset
	csi.clientSet = kubernetes.NewForConfigOrDie(csi.kubeConfig)
	csi.openshiftConfigClientSet = configclient.NewForConfigOrDie(rest.AddUserAgent(csi.kubeConfig, clientOperatorName))
}

func (csi *VsphereCSIDebugger) loginToVCenter(ctx context.Context) error {
	logger.SetLoggerLevel(logger.DevelopmentLogLevel)
	infra, err := csi.openshiftConfigClientSet.ConfigV1().Infrastructures().Get(ctx, "cluster", metav1.GetOptions{})
	if err != nil {
		debug("error getting infrastructure object: %v", infra)
		return err
	}

	debug("Creating vSphere connection")

	cloudConfig := infra.Spec.CloudConfig
	cloudConfigMap, err := csi.clientSet.CoreV1().ConfigMaps(cloudConfigNamespace).Get(ctx, cloudConfig.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get cloud config: %v", err)
	}

	csiConfigMap, err := csi.clientSet.CoreV1().ConfigMaps(csiDriverNamespace).Get(ctx, vsphereCsiConfig, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to load csi driver config: %v", err)
	}
	err = csi.writeCsiConfig(csiConfigMap)
	if err != nil {
		return fmt.Errorf("failed to write csi driver config: %v", err)
	}
	os.Setenv(EnvVSphereCSIConfig, csiConfigLocation)

	cfgString, ok := cloudConfigMap.Data[infra.Spec.CloudConfig.Key]
	if !ok {
		return fmt.Errorf("cloud config %s/%s does not contain key %q", cloudConfigNamespace, cloudConfig.Name, cloudConfig.Key)
	}

	var config vsphere.VSphereConfig
	err = gcfg.ReadStringInto(&config, cfgString)
	if err != nil {
		return err
	}
	vCenter := config.Workspace.VCenterIP

	secret, err := csi.clientSet.CoreV1().Secrets(csiDriverNamespace).Get(ctx, cloudCredSecretName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	userKey := vCenter + "." + "username"
	username, ok := secret.Data[userKey]
	if !ok {
		return fmt.Errorf("error parsing secret %q: key %q not found", cloudCredSecretName, userKey)
	}
	os.Setenv("VSPHERE_USER", string(username))

	config.Global.User = string(username)
	passwordKey := vCenter + "." + "password"
	password, ok := secret.Data[passwordKey]
	if !ok {
		return fmt.Errorf("error parsing secret %q: key %q not found", cloudCredSecretName, passwordKey)
	}
	os.Setenv("VSPHERE_PASSWORD", string(password))
	csi.config = &config
	return nil
}

func (csi *VsphereCSIDebugger) writeCsiConfig(configMap *corev1.ConfigMap) error {
	cfgString, ok := configMap.Data["cloud.conf"]
	if !ok {
		return fmt.Errorf("error writing csi configmap")
	}
	file, err := os.Create(csiConfigLocation)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(cfgString)
	if err != nil {
		return err
	}
	return nil
}

func (csi *VsphereCSIDebugger) getVolumePath(ctx context.Context, pvName string) (string, error) {
	pv, err := csi.clientSet.CoreV1().PersistentVolumes().Get(ctx, pvName, metav1.GetOptions{})
	if err != nil {
		debug("error finding pv %s: %v", pvName, err)
		return "", err
	}
	if pv.Spec.VsphereVolume != nil {
		return pv.Spec.VsphereVolume.VolumePath, nil
	} else {
		return "", fmt.Errorf("not a valid vsphere volume")
	}
}

func (csi *VsphereCSIDebugger) RegisterCNSDisk(ctx context.Context, pvName string) error {
	csi.createClients()

	err := csi.loginToVCenter(ctx)
	if err != nil {
		debug("failed to load config: %v", err)
		return err
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		debug("failed to generate uuid: %v", err)
		return err
	}
	re := regexp.MustCompile(`\[([^\[\]]*)\]`)
	volumePath, err := csi.getVolumePath(ctx, pvName)
	if err != nil {
		debug("failed to get volume path for pv: %v", err)
		return err
	}

	if !re.MatchString(volumePath) {
		return fmt.Errorf("failed to extract datastore name from in-tree volume path: %q", volumePath)
	}

	datastoreFullPath := re.FindAllString(volumePath, -1)[0]
	vmdkPath := strings.TrimSpace(strings.TrimPrefix(volumePath, datastoreFullPath))
	datastoreFullPath = strings.Trim(strings.Trim(datastoreFullPath, "["), "]")
	datastorePathSplit := strings.Split(datastoreFullPath, "/")
	datastoreName := datastorePathSplit[len(datastorePathSplit)-1]
	datacenters := csi.config.Workspace.Datacenter
	user := csi.config.Global.User
	host := csi.config.Workspace.VCenterIP

	// Get vCenter.
	vCenter, err := csi.GetVirtualCenter(ctx)
	if err != nil {
		debug("failed to get vCenter. err: %v", err)
		return err
	}
	manager, err := cnsvolume.GetManager(ctx, vCenter, nil, false, false, false, cnstypes.CnsClusterFlavorVanilla)
	if err != nil {
		debug("failed to get manager. err: %v", err)
		return err
	}

	datacenterPaths := make([]string, 0)
	if datacenters != "" {
		datacenterPaths = strings.Split(datacenters, ",")
	} else {
		// Get all datacenters from vCenter.
		dcs, err := vCenter.GetDatacenters(ctx)
		if err != nil {
			debug("failed to get datacenters from vCenter. err: %v", err)
			return err
		}
		for _, dc := range dcs {
			datacenterPaths = append(datacenterPaths, dc.InventoryPath)
		}
		debug("retrieved all datacenters %v from vCenter", datacenterPaths)
	}
	debug("found datacenter : %+v", datacenterPaths)
	var volumeInfo *cnsvolume.CnsVolumeInfo
	var storagePolicyID string
	var containerClusterArray []cnstypes.CnsContainerCluster
	containerCluster := cnsvsphere.GetContainerCluster("co8-dskwd", user, cnstypes.CnsClusterFlavorVanilla, "vanilla")

	containerClusterArray = append(containerClusterArray, containerCluster)
	createSpec := &cnstypes.CnsVolumeCreateSpec{
		Name:       uuid.String(),
		VolumeType: common.BlockVolumeType,
		Metadata: cnstypes.CnsVolumeMetadata{
			ContainerCluster:      containerCluster,
			ContainerClusterArray: containerClusterArray,
		},
	}
	if storagePolicyID != "" {
		profileSpec := &vim25types.VirtualMachineDefinedProfileSpec{
			ProfileId: storagePolicyID,
		}
		createSpec.Profile = append(createSpec.Profile, profileSpec)
	}
	for _, datacenter := range datacenterPaths {
		// Check vCenter API Version
		// Format:
		// https://<vc_ip>/folder/<vm_vmdk_path>?dcPath=<datacenter-path>&dsName=<datastoreName>
		backingDiskURLPath := "https://" + host + "/folder/" +
			vmdkPath + "?dcPath=" + url.PathEscape(datacenter) + "&dsName=" + url.PathEscape(datastoreName)
		bUseVslmAPIs, err := common.UseVslmAPIs(ctx, vCenter.Client.ServiceContent.About)
		if err != nil {
			return fmt.Errorf("Error while determining the correct APIs to use for vSphere version %q, Error= %+v", vCenter.Client.ServiceContent.About.ApiVersion, err)
		}
		if bUseVslmAPIs {
			backingObjectID, err := csi.RegisterDisk(ctx, backingDiskURLPath, volumePath)
			if err != nil {
				return fmt.Errorf("registration failed for volumePath: %v", volumePath)
			}
			createSpec.BackingObjectDetails = &cnstypes.CnsBlockBackingDetails{BackingDiskId: backingObjectID}
			debug("Registering volume: %q using backingDiskId :%q", volumePath, backingObjectID)
		} else {
			createSpec.BackingObjectDetails = &cnstypes.CnsBlockBackingDetails{BackingDiskUrlPath: backingDiskURLPath}
			debug("Registering volume: %q using backingDiskURLPath :%q", volumePath, backingDiskURLPath)
		}
		debug("vSphere CSI driver registering volume %q with create spec %+v", volumePath, spew.Sdump(createSpec))
		volumeInfo, _, err = manager.CreateVolume(ctx, createSpec, nil)
		if err != nil {
			debug("failed to register volume %q with createSpec: %v. error: %+v", volumePath, createSpec, err)
		} else {
			break
		}
	}
	if volumeInfo != nil {
		debug("Successfully registered volume %q as container volume with ID: %q", volumePath, volumeInfo.VolumeID.Id)
	} else {
		return fmt.Errorf("registration failed for volumeSpec: %v", volumePath)
	}
	debug("Successfully registered volume %q as container volume with ID: %q", volumePath, volumeInfo.VolumeID.Id)
	return nil
}

func (csi *VsphereCSIDebugger) GetVirtualCenter(ctx context.Context) (*cnsvsphere.VirtualCenter, error) {
	virtualCenterConfig, err := csi.GetVirtualCenterConfig(ctx)
	if err != nil {
		return nil, err
	}

	virtualCenter := &cnsvsphere.VirtualCenter{
		ClientMutex: &sync.Mutex{},
		Config:      virtualCenterConfig,
	}

	err = virtualCenter.Connect(ctx)
	if err != nil {
		return nil, err
	}
	csi.vCenter = virtualCenter
	return virtualCenter, nil
}

func (csi *VsphereCSIDebugger) GetVirtualCenterConfig(ctx context.Context) (*cnsvsphere.VirtualCenterConfig, error) {
	driverConfig, err := cnsconfig.GetConfig(ctx)
	if err != nil {
		debug("failed to get vsphere config: %v", err)
		return nil, err
	}
	virtualCenterConfigs, err := cnsvsphere.GetVirtualCenterConfigs(ctx, driverConfig)
	if err != nil {
		debug("failed to get virtual center configs: %v", err)
		return nil, err

	}
	for _, virtualCenterConfig := range virtualCenterConfigs {
		return virtualCenterConfig, nil
	}
	return nil, fmt.Errorf("no virtual center found")
}

func (csi *VsphereCSIDebugger) RegisterDisk(ctx context.Context, path string, name string) (string, error) {
	// Set up the VC connection.
	err := csi.vCenter.ConnectVslm(ctx)
	if err != nil {
		debug("ConnectVslm failed with err: %+v", err)
		return "", err
	}
	globalObjectManager := vslm.NewGlobalObjectManager(csi.vCenter.VslmClient)
	vStorageObject, err := globalObjectManager.RegisterDisk(ctx, path, name)
	if err != nil {
		alreadyExists, objectID := cnsvsphere.IsAlreadyExists(err)
		if alreadyExists {
			debug("vStorageObject: %q, already exists and registered as FCD, returning success", objectID)
			return objectID, nil
		}
		debug("failed to register virtual disk %q as first class disk with err: %v", path, err)
		return "", err
	}
	return vStorageObject.Config.Id.Id, nil
}
func debug(message string, args ...interface{}) {
	message = message + "\n"
	fmt.Printf(message, args...)
}
