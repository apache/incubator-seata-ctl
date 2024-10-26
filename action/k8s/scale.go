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

const (
	DefaultReplicas = 1
)

var ScaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "scale seata in k8s",
	Run: func(cmd *cobra.Command, args []string) {
		err := scale()
		if err != nil {
			tool.Logger.Errorf("scale err:%v", err)
		}
	},
}

func init() {
	ScaleCmd.PersistentFlags().StringVar(&Name, "name", "example-seataserver", "Seataserver name")
	ScaleCmd.PersistentFlags().Int32Var(&Replicas, "replicas", DefaultReplicas, "Replicas number")
	ScaleCmd.PersistentFlags().StringVar(&Namespace, "namespace", "default", "Namespace name")
}

func scale() error {
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
	if seataServer == nil {
		return fmt.Errorf("This seata server does not exits！" + Name)
	}
	seataServer.Object["spec"].(map[string]interface{})["replicas"] = Replicas

	_, err = client.Resource(gvr).Namespace(namespace).Update(context.TODO(), seataServer, metav1.UpdateOptions{})

	if err != nil {
		return err
	}
	tool.Logger.Infof("CR scale success，name: %s\n", Name)
	return nil
}
