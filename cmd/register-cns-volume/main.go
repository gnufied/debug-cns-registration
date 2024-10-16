package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/gnufied/debug-cns-registration/pkg/vsphere"
	"k8s.io/client-go/tools/clientcmd"
	klog "k8s.io/klog/v2"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// Some other error occurred, possibly indicating a problem
	return false
}

func main() {
	var kubeconfig *string
	pvName := flag.String("pv", "", "name of pv to migrate")
	klog.InitFlags(nil)
	flag.Parse()

	kubeConfigEnv := os.Getenv("KUBECONFIG")
	if kubeConfigEnv != "" && fileExists(kubeConfigEnv) {
		kubeconfig = &kubeConfigEnv
	}

	if pvName == nil || *pvName == "" {
		klog.Fatalf("Specify pv Name to migrate")
	}

	fmt.Printf("KubeConfig is: %s\n", *kubeconfig)

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		klog.Fatalf("error building kubeconfig: %v", err)
	}
	debugger, err := vsphere.NewVSphereCSIDebugger(config)
	if err != nil {
		fmt.Printf("error starting debugger: %v", err)
		return
	}
	err = debugger.RegisterCNSDisk(context.TODO(), *pvName)
	if err != nil {
		fmt.Printf("error registering disk: %v", err)
		return
	}
}
