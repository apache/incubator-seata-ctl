package k8s

import (
	"context"
	"fmt"
	"github.com/seata/seata-ctl/action/k8s/utils"
	"github.com/spf13/cobra"
	_ "gopkg.in/yaml.v3"
	_ "io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	_ "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	_ "k8s.io/apimachinery/pkg/runtime/schema"
	_ "k8s.io/client-go/applyconfigurations/meta/v1"
	"log"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Kubernetes CRD controller",
	Run: func(cmd *cobra.Command, args []string) {
		err := DeployCRD()
		if err != nil {
			log.Fatal(err)
		}
		// err = DeployController()
		// if err != nil {
		// 	log.Fatal(err)
		// }
	},
}

const CreateCrdPath = "/apis/apiextensions.k8s.io/v1/customresourcedefinitions"
const FilePath = "seata.yaml"

// DeployCRD deploys the custom resource definition.
func DeployCRD() error {
	res, err := utils.CreateRequest(CreateCrdPath, FilePath)
	fmt.Println(res)
	if err != nil {
		return err
	}
	return nil
}

// DeployController deploys the controller for the custom resource.
func DeployController() error {
	// Get Kubernetes client
	clientset, err := utils.GetClient()
	if err != nil {
		return fmt.Errorf("error getting clientset: %v", err)
	}
	// Create Deployment object
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "seata-k8s-controller-manager",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func(i int32) *int32 { return &i }(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "seata-k8s-controller-manager",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "seata-k8s-controller-manager",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "seata-k8s-controller-manager",
							Image: "apache/seata-controller:latest",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = clientset.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Error creating deployment: %s", err.Error())
	}
	fmt.Println("Deployment created successfully")
	return err
}
