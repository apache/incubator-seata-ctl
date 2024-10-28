package k8s

import (
	"context"
	"github.com/seata/seata-ctl/action/k8s/utils"
	"github.com/seata/seata-ctl/tool"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var UnDeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "undeploy seata in k8s",
	Run: func(cmd *cobra.Command, args []string) {
		err := undeploy()
		if err != nil {
			tool.Logger.Errorf("undeploy error: %v", err)
		}
	},
}

func init() {
	UnDeployCmd.PersistentFlags().StringVar(&Name, "name", DefaultCRName, "Seataserver name")
	UnDeployCmd.PersistentFlags().StringVar(&Namespace, "namespace", DefaultNamespace, "Namespace name")
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

	tool.Logger.Infof("CR delete success，name: %s\n", Name)

	return nil
}
