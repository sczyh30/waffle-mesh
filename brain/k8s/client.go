package k8s

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
)

func CreateInClusterInterface() (*rest.Config, kubernetes.Interface, error) {
	restConfig, err1 := rest.InClusterConfig()
	if err1 == nil {
		client, err2 := kubernetes.NewForConfig(restConfig)
		if err2 == nil {
			return restConfig, client, nil
		}
		return nil, nil, err2
	}
	return nil, nil, err1
}
