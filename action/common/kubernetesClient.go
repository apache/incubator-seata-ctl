package common

import (
	"github.com/seata/seata-ctl/model"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

// GetKubeConfigPath 根据配置文件中的内容，获取kubeConfigPath
func GetKubeConfigPath() (string, error) {
	//获取配置文件
	file, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Failed to read config.yml: %v", err)
	}
	var config model.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	//获取context name
	contextName := config.Context.Kubernetes
	var contextPath string
	//根据context，获取匹配的kubeconfig位置
	for _, cluster := range config.Kubernetes.Cluster {
		if cluster.Name == contextName {
			contextPath = cluster.KubeConfigPath
		}
	}
	//如果没有找到，返回错误
	if contextPath == "" {
		log.Fatalf("Failed to find context in config.yml")
	}
	return contextPath, err
}

// GetClient 根据指定的 kubeconfig 文件路径创建并返回 Kubernetes 客户端
func GetClient() (*kubernetes.Clientset, error) {
	// 使用 client 加载指定的 kubeconfig 文件
	kubeconfigPath, err := GetKubeConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// 创建一个 Kubernetes 客户端
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetDynamicClient 根据指定的 kubeconfig 文件路径创建并返回 Kubernetes 动态客户端
func GetDynamicClient() (*dynamic.DynamicClient, error) {
	// 使用 client 加载指定的 kubeconfig 文件
	kubeconfigPath, err := GetKubeConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// 创建动态客户端
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("无法创建动态客户端: %v", err)
	}
	return dynamicClient, nil
}
