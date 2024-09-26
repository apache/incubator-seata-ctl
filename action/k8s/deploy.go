package k8s

import (
	"context"
	"fmt"
	"github.com/seata/seata-ctl/action/k8s/utils"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy seata in k8s",
	Run: func(cmd *cobra.Command, args []string) {
		deploy()
	},
}

var Name string
var ServiceName string
var Replicas int32
var Image string
var Namespace string

func init() {
	DeployCmd.PersistentFlags().StringVar(&Name, "name", "example-seataserver", "Seataserver name")
	DeployCmd.PersistentFlags().StringVar(&ServiceName, "service", "seata-service", "Headless Service name")
	DeployCmd.PersistentFlags().Int32Var(&Replicas, "replicas", 1, "Replicas number")
	DeployCmd.PersistentFlags().StringVar(&Image, "image", "seata-server", "Image tag")
	DeployCmd.PersistentFlags().StringVar(&Namespace, "namespace", "default", "Namespace name")
}

func deploy() error {
	//获取动态kubeclient
	client, err := utils.GetDynamicClient()
	if err != nil {
		return err
	}
	// 获取命名空间
	namespace := Namespace

	// 定义 Custom Resource 的 GroupVersionResource
	gvr := schema.GroupVersionResource{
		Group:    "operator.seata.apache.org",
		Version:  "v1alpha1",
		Resource: "seataservers",
	}

	var seataServer *unstructured.Unstructured
	seataServer, err = client.Resource(gvr).Namespace(namespace).Get(context.TODO(), Name, metav1.GetOptions{})
	if seataServer != nil {
		fmt.Println("This seata server already exits！")
		return nil
	}
	// 创建 Custom Resource 对象
	seataServer = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "operator.seata.apache.org/v1alpha1",
			"kind":       "SeataServer",
			"metadata": map[string]interface{}{
				"name":      Name,
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"containerName": Image,
				"image":         "seataio/seata-server:latest",
				"replicas":      Replicas,
				"serviceName":   ServiceName,
				//"env": map[string]interface{}{
				//	"SEATA_ENV":       "prod",
				//	"SEATA_LOG_LEVEL": "info",
				//},
				"ports": map[string]interface{}{
					"consolePort": 7091,
					"raftPort":    9091,
					"servicePort": 8091,
				},
				"resources": map[string]interface{}{
					"limits": map[string]interface{}{
						"cpu":    "500m",
						"memory": "1Gi",
					},
					"requests": map[string]interface{}{
						"cpu":    "250m",
						"memory": "512Mi",
					},
				},
				"store": map[string]interface{}{
					"resources": map[string]interface{}{
						"limits": map[string]interface{}{
							"cpu":    "500m",
							"memory": "1Gi",
						},
						"requests": map[string]interface{}{
							"cpu":    "250m",
							"memory": "512Mi",
						},
					},
				},
			},
		},
	}
	// 尝试创建 Custom Resource
	_, err = client.Resource(gvr).Namespace(namespace).Create(context.TODO(), seataServer, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("CR install success，name: %s\n", Name)
	return nil
}
