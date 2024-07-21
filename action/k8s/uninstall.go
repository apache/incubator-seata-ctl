package k8s

import (
	"context"
	"github.com/seata/seata-ctl/action/common"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var UnDeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "undeploy seata in kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		undeploy()
	},
}

func undeploy() error {
	//获取动态kubeclient
	client, err := common.GetDynamicClient()
	if err != nil {
		return err
	}

	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}

	//删除 CRD
	err = client.Resource(gvr).Delete(context.TODO(), "name", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
