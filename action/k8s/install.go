package k8s

import (
	"context"
	"fmt"
	"github.com/seata/seata-ctl/action/common"
	"github.com/spf13/cobra"
	_ "gopkg.in/yaml.v3"
	_ "io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	_ "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	_ "k8s.io/apimachinery/pkg/runtime/schema"
	_ "k8s.io/client-go/applyconfigurations/meta/v1"
	"log"
	"os/exec"
	"strings"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install k8s crd controller",
	Run: func(cmd *cobra.Command, args []string) {
		//err := InstallCRD()
		//if err != nil {
		//	log.Fatal(err)
		//}
		err := DeploymentController()
		if err != nil {
			log.Fatal(err)
		}
		output, err := executeCommand("kubectl apply -f operator.seata.apache.org_seataservers.yaml")
		if err != nil {
			log.Println(output)
		}
	},
}

func executeCommand(input string) (string, error) {
	// 将用户输入的命令分割成命令和参数
	args := strings.Fields(input)
	if len(args) == 0 {
		return "", fmt.Errorf("no command provided")
	}

	// 使用 exec.Command 执行命令
	cmd := exec.Command(args[0], args[1:]...)

	// 获取命令的输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// InstallCRD 部署CRD
func InstallCRD() error {
	//获取动态kubeclient
	client, err := common.GetDynamicClient()
	if err != nil {
		return err
	}

	// 定义 CRD 对象
	crd := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apiextensions.k8s.io/v1",
			"kind":       "CustomResourceDefinition",
			"metadata": map[string]interface{}{
				"name": "seataservers.operator.seata.apache.org",
			},
			"spec": map[string]interface{}{
				"group": "operator.seata.apache.org",
				"names": map[string]interface{}{
					"kind":     "SeataServer",
					"listKind": "SeataServerList",
					"plural":   "seataservers",
					"singular": "seataserver",
				},
				"scope": "Namespaced",
				"versions": []map[string]interface{}{
					{
						"name":    "v1alpha1",
						"served":  true,
						"storage": true,
						"schema": map[string]interface{}{
							"openAPIV3Schema": map[string]interface{}{
								"description": "SeataServer is the Schema for the seataservers API",
								"type":        "object",
								"properties": map[string]interface{}{
									"apiVersion": map[string]interface{}{
										"description": "APIVersion defines the versioned schema of this representation",
										"type":        "string",
									},
									"kind": map[string]interface{}{
										"description": "Kind is a string value representing the REST resource this",
										"type":        "string",
									},
									"metadata": map[string]interface{}{
										"type": "object",
									},
									"spec": map[string]interface{}{
										"description": "SeataServerSpec defines the desired state of SeataServer",
										"type":        "object",
										"properties": map[string]interface{}{
											"containerName": map[string]interface{}{
												"default": "seata-server",
												"type":    "string",
											},
											"env": map[string]interface{}{
												"additionalProperties": map[string]interface{}{
													"type": "string",
												},
												"type": "object",
											},
											"image": map[string]interface{}{
												"type": "string",
											},
											"ports": map[string]interface{}{
												"properties": map[string]interface{}{
													"consolePort": map[string]interface{}{
														"default": 7091,
														"format":  "int32",
														"type":    "integer",
													},
													"raftPort": map[string]interface{}{
														"default": 9091,
														"format":  "int32",
														"type":    "integer",
													},
													"servicePort": map[string]interface{}{
														"default": 8091,
														"format":  "int32",
														"type":    "integer",
													},
												},
												"type": "object",
											},
											"replicas": map[string]interface{}{
												"default": 1,
												"format":  "int32",
												"minimum": 1,
												"type":    "integer",
											},
											//"resources": map[string]interface{}{
											//	"description": "ResourceRequirements describes the compute resource requirements.",
											//	"properties": map[string]interface{}{
											//		"claims": map[string]interface{}{
											//			"description": "Claims lists the names of resources, defined in spec.resourceClaims, that are used by this container.This is an alpha field and requires enabling the DynamicResourceAllocation feature gate. This field is immutable. It can only be set for containers.",
											//			"items": map[string]interface{}{
											//				"description": "ResourceClaim references one entry in PodSpec.ResourceClaims.",
											//				"properties": map[string]interface{}{
											//					"name": map[string]interface{}{
											//						"description": "Name must match the name of one entry in pod.spec.resourceClaims of the Pod where this field is used. It makes that resource available inside a container.",
											//						"type":        "string",
											//					},
											//				},
											//				"required": "name",
											//				"type": "object",
											//			},
											//			"type":"array",
											//			"x-k8s-list-map-keys": "name",
											//			"x-k8s-list-type":     "map",
											//		},
											//		"limits": map[string]interface{}{
											//			"additionalProperties": map[string]interface{}{
											//				"anyOf": map[string]interface{}{
											//					//"type": "integer",
											//					"type": "string",
											//				},
											//				"pattern":                    "^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$",
											//				"x-k8s-int-or-string": "true",
											//			},
											//			"description": "Limits describes the maximum amount of compute resources",
											//			"type":        "object",
											//		},
											//		"requests": map[string]interface{}{
											//			"additionalProperties": map[string]interface{}{
											//				"anyOf": map[string]interface{}{
											//					//"type": "integer",
											//					"type": "string",
											//				},
											//				"pattern":                    "^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$",
											//				"x-k8s-int-or-string": "true",
											//			},
											//			"description": "Requests describes the minimum amount of compute resources required. If Requests is omitted for a container,it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits.More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers.",
											//			"type":        "object",
											//		},
											//	},
											//	"type": "object",
											//},
											"serviceName": map[string]interface{}{
												"default": "seata-server",
												"type":    "string",
											},
											//"store": map[string]interface{}{
											//	"properties": map[string]interface{}{
											//		"resources": map[string]interface{}{
											//			"description": "ResourceRequirements describes the compute resource requirements.",
											//			"properties": map[string]interface{}{
											//				"claims": map[string]interface{}{
											//					"description": "Claims lists the names of resources, defined",
											//					"items": map[string]interface{}{
											//						"description": "ResourceClaim references one entry in PodSpec.ResourceClaims.",
											//						"properties": map[string]interface{}{
											//							"name": map[string]interface{}{
											//								"description": "Name must match the name of one entry in pod.spec.resourceClaims",
											//								"type":        "string",
											//							},
											//						},
											//						"required": "name",
											//						"type":     "object",
											//					},
											//					"type":                       "array",
											//					"x-k8s-list-map-keys": "name",
											//					"x-k8s-list-type":     "map",
											//				},
											//				"limits": map[string]interface{}{
											//					"additionalProperties": map[string]interface{}{
											//						"anyOf": map[string]interface{}{
											//							"type": "integer",
											//							//"type":"string",
											//						},
											//						"pattern":                    "^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$",
											//						"x-k8s-int-or-string": "true",
											//					},
											//					"description": "Limits describes the maximum amount of compute resources",
											//					"type":        "object",
											//				},
											//				"requests": map[string]interface{}{
											//					"additionalProperties": map[string]interface{}{
											//						"anyOf": map[string]interface{}{
											//							"type": "integer",
											//							//"type":"string",
											//						},
											//						"pattern":                    "^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$",
											//						"x-k8s-int-or-string": "true",
											//					},
											//					"description": "Requests describes the minimum amount of compute",
											//					"type":        "object",
											//				},
											//				"type": "object",
											//			},
											//			"required": "resources",
											//			"type":     "object",
											//		},
										},
									},
									"status": map[string]interface{}{
										"description": "SeataServerStatus defines the observed state of SeataServer",
										"properties": map[string]interface{}{
											"readyReplicas": map[string]interface{}{
												"format": "int32",
												"type":   "integer",
											},
											"replicas": map[string]interface{}{
												"format": "int32",
												"type":   "integer",
											},
											"synchronized": map[string]interface{}{
												"type": "boolean",
											},
										},
										//"required": "replicas",
										"type": "object",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}

	//创建 CRD
	result, err := client.Resource(gvr).Create(context.TODO(), crd, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			log.Printf("CRD 已存在: %v", err)
		} else {
			log.Fatalf("创建 CRD 失败: %v", err)
		}
	}

	fmt.Printf("CRD 创建成功，名称: %s\n", result.GetName())
	return nil
}

// DeploymentController 部署Seata-controller
func DeploymentController() error {
	//获取kubeclient
	clientset, err := common.GetClient()
	//if err != nil {
	//	return err
	//}
	// 创建Deployment对象
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

	// 使用客户端创建Deployment
	_, err = clientset.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Error creating deployment: %s", err.Error())
	}
	fmt.Println("Deployment created successfully")
	return err
}
