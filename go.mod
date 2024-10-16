module github.com/gnufied/debug-cns-registration

go 1.22.0

toolchain go1.23.1

require (
	github.com/container-storage-interface/spec v1.9.0
	github.com/davecgh/go-spew v1.1.1
	github.com/google/uuid v1.6.0
	github.com/openshift/client-go v0.0.0-00010101000000-000000000000
	github.com/vmware/govmomi v0.37.3
	gopkg.in/gcfg.v1 v1.2.3
	k8s.io/apimachinery v0.27.10
	k8s.io/client-go v0.27.10
	k8s.io/klog/v2 v2.100.1
	k8s.io/legacy-cloud-providers v0.0.0
	sigs.k8s.io/vsphere-csi-driver/v3 v3.3.1
)

require (
	github.com/Microsoft/go-winio v0.4.17 // indirect
	github.com/akutz/gofsutil v0.1.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.1 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kubernetes-csi/csi-proxy/client v1.1.3 // indirect
	github.com/kubernetes-csi/external-snapshotter/client/v6 v6.1.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/selinux v1.10.0 // indirect
	github.com/openshift/api v0.0.0-20230807121159-a81c3efc8824 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.18.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vmware-tanzu/vm-operator/api v1.8.2 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/oauth2 v0.13.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/term v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.3.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231127180814-3a041ad873d4 // indirect
	google.golang.org/grpc v1.60.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/warnings.v0 v0.1.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.27.10 // indirect
	k8s.io/apiextensions-apiserver v0.27.10 // indirect
	k8s.io/apiserver v0.27.10 // indirect
	k8s.io/cloud-provider v0.26.10 // indirect
	k8s.io/component-base v0.27.10 // indirect
	k8s.io/component-helpers v0.27.10 // indirect
	k8s.io/kube-openapi v0.0.0-20230501164219-8b0f38b5fd1f // indirect
	k8s.io/kubernetes v1.27.10 // indirect
	k8s.io/mount-utils v0.27.10 // indirect
	k8s.io/sample-controller v0.27.10 // indirect
	k8s.io/utils v0.0.0-20230505201702-9f6742963106 // indirect
	sigs.k8s.io/controller-runtime v0.15.3 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20230807132528-be5346fb33cb
	github.com/openshift/library-go => github.com/openshift/library-go v0.0.0-20231213084759-840298df1eee
	k8s.io/legacy-cloud-providers => github.com/kubernetes/legacy-cloud-providers v0.0.0-20240512054615-314470656d51

)
