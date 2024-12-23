package kubernetes

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// get kubernetes config

func GetKubeConfig() (*kubernetes.Clientset, error) {
	config := rest.Config{}
	cfg := &config
	var err error

	cfg, err = rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to create in cluster config: %v", err)
		return nil, err
	}

	kubeClientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
		return nil, err
	}

	return kubeClientSet, nil
}
