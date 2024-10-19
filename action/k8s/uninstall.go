package k8s

import (
	"context"
	"fmt"
	"github.com/seata/seata-ctl/action/k8s/utils"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
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

	// Assume client and gvr have already been defined
	err = client.Resource(gvr).Delete(context.TODO(), CRDname, metav1.DeleteOptions{})
	if err != nil {
		// Check if the error is a "not found" error
		if errors.IsNotFound(err) {
			// The resource does not exist, output a message instead of returning an error
			fmt.Printf("CRD %s does not exist, no action taken.\n", CRDname)
		} else {
			// For other errors, log the error and exit the program
			log.Fatalf("delete CRD failed: %v", err)
		}
	} else {
		// Successfully deleted the resource
		fmt.Printf("CRD %s deleted successfully.\n", CRDname)
	}

	return nil
}

func UnDeploymentController() error {
	clientset, err := utils.GetClient()
	if err != nil {
		return err
	}
	// Assume clientset has already been defined
	err = clientset.AppsV1().Deployments(Namespace).Delete(context.TODO(), Deployname, metav1.DeleteOptions{})
	if err != nil {
		// Check if the error is a "not found" error
		if errors.IsNotFound(err) {
			// The deployment does not exist, output a message instead of returning an error
			fmt.Printf("Deployment 'seata-k8s-controller-manager' does not exist in namespace '%s', no action taken.\n", Namespace)
		} else {
			// For other errors, log the error and exit the program
			log.Fatalf("Error deleting deployment: %s", err.Error())
		}
	} else {
		// Successfully deleted the deployment
		fmt.Printf("Deployment 'seata-k8s-controller-manager' deleted successfully from namespace '%s'.\n", Namespace)
	}
	return nil
}
