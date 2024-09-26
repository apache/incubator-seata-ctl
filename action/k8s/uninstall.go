package k8s

import (
	"context"
	"fmt"
	"github.com/seata/seata-ctl/action/k8s/utils"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
)

var UnInstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall seata in k8s",
	Run: func(cmd *cobra.Command, args []string) {
		err := UninstallCRD()
		if err != nil {
			log.Fatal(err)
		}
		err = UnDeploymentController()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// UninstallCRD 删除CRD定义
func UninstallCRD() error {
	//获取动态kubeclient
	client, err := utils.GetDynamicClient()
	if err != nil {
		return err
	}
	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}
	//删除 CRD
	err = client.Resource(gvr).Delete(context.TODO(), "seataservers.operator.seata.apache.org", metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("删除CRD 失败: %v", err)
	}
	fmt.Printf("delete CRD success，名称: %s\n", "seataservers.operator.seata.apache.org")
	return nil
}

// UnDeploymentController 删除Seata-controller
func UnDeploymentController() error {
	clientset, err := utils.GetClient()
	if err != nil {
		return err
	}
	err = clientset.AppsV1().Deployments("default").Delete(context.TODO(), "seata-k8s-controller-manager", metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("Error creating deployment: %s", err.Error())
	}
	fmt.Println("Deployment created successfully")
	return err
}
