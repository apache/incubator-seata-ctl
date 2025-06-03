package k8s

import (
	"context"
	"fmt"

	"github.com/seata/seata-ctl/action/k8s/utils"
	"github.com/seata/seata-ctl/tool"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Kubernetes CRD controller",
	Run: func(_ *cobra.Command, _ []string) {
		err := DeployCRD()
		if err != nil {
			tool.Logger.Errorf("install CRD err: %v", err)
		}
		err = DeployController()
		if err != nil {
			tool.Logger.Errorf("install Controller err: %v", err)
		}
	},
}

func init() {
	InstallCmd.PersistentFlags().StringVar(&Namespace, "namespace", DefaultNamespace, "Namespace name")
	InstallCmd.PersistentFlags().StringVar(&ControllerImage, "image", DefaultControllerImage, "Namespace name")
	InstallCmd.PersistentFlags().StringVar(&DeployName, "name", DefaultDeployName, "Deployment name")
}

// DeployCRD deploys the custom resource definition.
func DeployCRD() error {
	_, err := utils.CreateRequest(CreateCrdPath, FilePath)
	if err != nil {
		return err
	}
	return nil
}

// DeployController deploys the controller for the custom resource.
func DeployController() error {
	// Get Kubernetes client
	client, err := utils.GetClient()
	if err != nil {
		return fmt.Errorf("get client err: %v", err)
	}

	// Define the Deployment name and namespace
	deploymentName := DeployName
	namespace := Namespace

	// Check if the Deployment already exists
	_, err = client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err == nil {
		// If the Deployment exists, output a message and return
		return fmt.Errorf("deployment '%s' already exists in the '%s' namespace", deploymentName, Namespace)
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
	_, err = client.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating deployment: %v", err)
	}
	tool.Logger.Infof("Deployment created successfully")
	return nil
}
