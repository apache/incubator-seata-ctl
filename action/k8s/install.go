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
	"k8s.io/apimachinery/pkg/api/errors"
	_ "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	_ "k8s.io/apimachinery/pkg/runtime/schema"
	_ "k8s.io/client-go/applyconfigurations/meta/v1"
	"log"
)

const (
	CreateCrdPath = "/apis/apiextensions.k8s.io/v1/customresourcedefinitions"
	FilePath      = "seata.yaml"

	CRDname    = "seataservers.operator.seata.apache.org"
	Deployname = "seata-k8s-controller-manager"

	DefaultControllerImage = "apache/seata-controller:latest"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Kubernetes CRD controller",
	Run: func(cmd *cobra.Command, args []string) {
		err := DeployCRD()
		if err != nil {
			log.Fatal(err)
		}
		err = DeployController()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var ControllerImage string

func init() {
	InstallCmd.PersistentFlags().StringVar(&Namespace, "namespace", "default", "Namespace name")
	InstallCmd.PersistentFlags().StringVar(&ControllerImage, "image", DefaultControllerImage, "Namespace name")
}

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

	// Define the Deployment name and namespace
	deploymentName := Deployname
	namespace := Namespace

	// Check if the Deployment already exists
	_, err = clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err == nil {
		// If the Deployment exists, output a message and return
		fmt.Printf("Deployment '%s' already exists in the '%s' namespace\n", deploymentName, Namespace)
		return nil
	} else if !errors.IsNotFound(err) {
		// If there is an error other than "not found", return it
		return fmt.Errorf("error checking for existing deployment: %v", err)
	}

	// Create Deployment object if it does not exist
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func(i int32) *int32 { return &i }(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  deploymentName,
							Image: ControllerImage,
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

	// Create the Deployment
	_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating deployment: %v", err)
	}
	fmt.Println("Deployment created successfully")
	return nil
}
