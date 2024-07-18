package common

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClient 根据指定的 kubeconfig 文件路径创建并返回 Kubernetes 客户端
func GetClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	// 使用 clientcmd 加载指定的 kubeconfig 文件
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
