package k8s

import (
	"context"
	"fmt"
	"github.com/seata/seata-ctl/action/k8s/utils"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var UnDeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "deploy seata in k8s",
	Run: func(cmd *cobra.Command, args []string) {
		err := undeploy()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	UnDeployCmd.PersistentFlags().StringVar(&Name, "name", "example-seataserver", "Seataserver name")
	UnDeployCmd.PersistentFlags().StringVar(&Namespace, "namespace", "default", "Namespace name")
}

func undeploy() error {
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

	err = client.Resource(gvr).Namespace(namespace).Delete(context.TODO(), Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("CR delete success，name: %s\n", Name)
	return nil
}
