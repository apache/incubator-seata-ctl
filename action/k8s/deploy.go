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

const (
	DefaultServerImage = "apache/seata-server:latest"
	ServiceName        = "seata-server-cluster"
	RequestStorage     = "1Gi"
	LimitStorage       = "1Gi"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy seata in k8s",
	Run: func(cmd *cobra.Command, args []string) {
		err := deploy()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var (
	Name      string
	Replicas  int32
	Namespace string
	Image     string
)

func init() {
	DeployCmd.PersistentFlags().StringVar(&Name, "name", "example-seataserver", "Seataserver name")
	DeployCmd.PersistentFlags().Int32Var(&Replicas, "replicas", 1, "Replicas number")
	DeployCmd.PersistentFlags().StringVar(&Namespace, "namespace", "default", "Namespace name")
	DeployCmd.PersistentFlags().StringVar(&Image, "image", DefaultServerImage, "Seata server image")
}

func deploy() error {
	client, err := utils.GetDynamicClient()
	if err != nil {
		return err
	}
	namespace := Namespace
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
	seataServer = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "operator.seata.apache.org/v1alpha1",
			"kind":       "SeataServer",
			"metadata": map[string]interface{}{
				"name":      Name,
				"namespace": Namespace,
			},
			"spec": map[string]interface{}{
				"serviceName": ServiceName,
				"replicas":    Replicas,
				"image":       Image,
				"store": map[string]interface{}{
					"resources": map[string]interface{}{
						"requests": map[string]interface{}{
							"storage": RequestStorage,
						},
						"limits": map[string]interface{}{
							"storage": LimitStorage,
						},
					},
				},
			},
		},
	}
	_, err = client.Resource(gvr).Namespace(namespace).Create(context.TODO(), seataServer, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("CR install success，name: %s\n", Name)
	return nil
}
