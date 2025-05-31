package utils

import (
	"fmt"
	"os"

	"github.com/seata/seata-ctl/model"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const ConfigFileName = "config.yml"

// GetKubeConfigPath retrieves the kubeConfigPath based on the contents of the config file.
func GetKubeConfigPath() (string, error) {
	// Read the configuration file
	file, err := os.ReadFile(ConfigFileName)
	if err != nil {
		return "", fmt.Errorf("failed to read config.yml:" + err.Error())
	}
	var config model.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return "", fmt.Errorf("unmarshal failed" + err.Error())
	}
	// Retrieve the context name
	contextName := config.Context.Kubernetes
	var contextPath string
	// Find the matching KubeConfig path based on the context
	for _, cluster := range config.Kubernetes.Cluster {
		if cluster.Name == contextName {
			contextPath = cluster.KubeConfigPath
		}
	}
	// If no matching context is found, return an error
	if contextPath == "" {
		return "", fmt.Errorf("failed to find context in config.yml")
	}
	return contextPath, err
}

// GetClient creates and returns a Kubernetes client based on the specified KubeConfigPath file path.
func GetClient() (*kubernetes.Clientset, error) {
	// Load the KubeConfig file using the client
	KubeConfigPath, err := GetKubeConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", KubeConfigPath)
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

// GetDynamicClient creates and returns a Kubernetes dynamic client based on the specified KubeConfigPath file path.
func GetDynamicClient() (*dynamic.DynamicClient, error) {
	// Load the KubeConfigPath file using the client
	KubeConfigPath, err := GetKubeConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", KubeConfigPath)
	if err != nil {
		return nil, err
	}

	// Create a dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client" + err.Error())
	}
	return dynamicClient, nil
}
