package utils

import (
	"github.com/seata/seata-ctl/model"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

// GetKubeConfigPath retrieves the kubeConfigPath based on the contents of the config file.
func GetKubeConfigPath() (string, error) {
	// Read the configuration file
	file, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Failed to read config.yml: %v", err)
	}
	var config model.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	// Retrieve the context name
	contextName := config.Context.Kubernetes
	var contextPath string
	// Find the matching kubeconfig path based on the context
	for _, cluster := range config.Kubernetes.Cluster {
		if cluster.Name == contextName {
			contextPath = cluster.KubeConfigPath
		}
	}
	// If no matching context is found, return an error
	if contextPath == "" {
		log.Fatalf("Failed to find context in config.yml")
	}
	return contextPath, err
}

// GetClient creates and returns a Kubernetes client based on the specified kubeconfig file path.
func GetClient() (*kubernetes.Clientset, error) {
	// Load the kubeconfig file using the client
	kubeconfigPath, err := GetKubeConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// Create a Kubernetes client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetDynamicClient creates and returns a Kubernetes dynamic client based on the specified kubeconfig file path.
func GetDynamicClient() (*dynamic.DynamicClient, error) {
	// Load the kubeconfig file using the client
	kubeconfigPath, err := GetKubeConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// Create a dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create dynamic client: %v", err)
	}
	return dynamicClient, nil
}
