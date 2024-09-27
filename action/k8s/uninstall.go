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

const CRDname = "seataservers.operator.seata.apache.org"

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

func init() {
	UnInstallCmd.PersistentFlags().StringVar(&Namespace, "namespace", "default", "Namespace name")
}

func UninstallCRD() error {
	client, err := utils.GetDynamicClient()
	if err != nil {
		return err
	}

	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}

	err = client.Resource(gvr).Delete(context.TODO(), CRDname, metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("delete CRD failed: %v", err)
	}
	fmt.Printf("delete CRD success，name: %s\n", "seataservers.operator.seata.apache.org")
	return nil
}

func UnDeploymentController() error {
	clientset, err := utils.GetClient()
	if err != nil {
		return err
	}
	err = clientset.AppsV1().Deployments(Namespace).Delete(context.TODO(), "seata-k8s-controller-manager", metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("Error creating deployment: %s", err.Error())
	}
	fmt.Println("Deployment created successfully")
	return err
}
