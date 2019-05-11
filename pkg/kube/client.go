package kube

import (
	"os"

	"github.com/golang/glog"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type kubeClient struct {
	*k8s.Clientset
}

// TODO
// modify configfile path later
func newKubeClient() *kubeClient {
	var cfg *rest.Config
	var err error
	// cPath := "/Users/tommenx/.kube/config"
	cPath := "/root/.kube/config"
	if cPath != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", cPath)
		if err != nil {
			glog.Errorf("Failed to get cluster config with error: %v\n", err)
			os.Exit(1)
		}
	} else {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			glog.Errorf("Failed to get cluster config with error: %v\n", err)
			os.Exit(1)
		}
	}
	client, err := k8s.NewForConfig(cfg)
	if err != nil {
		glog.Errorf("Failed to create client with error: %v\n", err)
		os.Exit(1)
	}
	return &kubeClient{client}
}
