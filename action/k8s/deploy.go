package k8s

import (
	"context"
	"fmt"

	"github.com/seata/seata-ctl/action/k8s/utils"
	"github.com/seata/seata-ctl/tool"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy seata in k8s",
	Run: func(_ *cobra.Command, _ []string) {
		err := deploy()
		if err != nil {
			tool.Logger.Errorf("deploy err:%v", err)
		}
	},
}

func init() {
	DeployCmd.PersistentFlags().StringVar(&Name, "name", DefaultCRName, "Seataserver name")
	DeployCmd.PersistentFlags().Int32Var(&Replicas, "replicas", 1, "Replicas number")
	DeployCmd.PersistentFlags().StringVar(&Namespace, "namespace", DefaultNamespace, "Namespace name")
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
	if err != nil {
		return err
	}
	if seataServer != nil {
		return fmt.Errorf("seata server already exist! name:" + Name)
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
	tool.Logger.Infof("CR install success，name: %s\n", Name)
	return nil
}
