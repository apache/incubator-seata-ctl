package k8s

import (
	"context"
	"fmt"
	"github.com/seata/seata-ctl/action/common"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install k8s crd controller",
	Run: func(cmd *cobra.Command, args []string) {
		err := install()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// 测试client 连接方法：输出所有的pod
func test() error {
	//获取kubeclient
	clientset, err := common.GetClient()
	if err != nil {
		return err
	}
	//部署crd资源
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %v", err)
	}

	for _, pod := range pods.Items {
		fmt.Printf("Namespace: %s, Name: %s\n", pod.Namespace, pod.Name)
	}
	return nil
}

func install() error {
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
				"name": "crontabs.stable.com",
			},
			"spec": map[string]interface{}{
				"group": "stable.com",
				"scope": "Namespaced",
				"names": map[string]interface{}{
					"plural":     "crontabs",
					"singular":   "crontab",
					"kind":       "CronTab",
					"shortNames": []string{"ct"},
				},
				"versions": []map[string]interface{}{
					{
						"name":    "v1",
						"served":  true,
						"storage": true,
						"schema": map[string]interface{}{
							"openAPIV3Schema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"spec": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"cronSpec": map[string]interface{}{
												"type":    "string",
												"pattern": "^((\\d+|\\*)(\\/\\d+)?){4}$",
											},
											"image": map[string]interface{}{
												"type": "string",
											},
											"replicas": map[string]interface{}{
												"type": "integer",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	//crd := &unstructured.Unstructured{
	//	Object: map[string]interface{}{
	//		"apiVersion": "apiextensions.k8s.io/v1",
	//		"kind":       "CustomResourceDefinition",
	//		"metadata": map[string]interface{}{
	//			"name": "seataservers.operator.seata.apache.org",
	//			"annotations": map[string]string{
	//				"controller-gen.kubebuilder.io/version": "v0.13.0",
	//			},
	//		},
	//		"spec": map[string]interface{}{
	//			"group": "operator.seata.apache.org",
	//			"scope": "Namespaced",
	//			"names": map[string]interface{}{
	//				"kind":     "SeataServer",
	//				"listKind": "SeataServerList",
	//				"plural":   "seataservers",
	//				"singular": "seataserver",
	//			},
	//			"versions": []interface{}{
	//				map[string]interface{}{
	//					"name": "v1alpha1",
	//					"schema": map[string]interface{}{
	//						"openAPIV3Schema": map[string]interface{}{
	//							"description": "SeataServer is the Schema for the seataservers API",
	//							"properties": map[string]interface{}{
	//								"apiVersion": map[string]interface{}{
	//									"description": "APIVersion defines the versioned schema of this representation\nof an object. Servers should convert recognized schemas to the latest\ninternal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
	//									"type":        "string",
	//								},
	//								"kind": map[string]interface{}{
	//									"description": "Kind is a string value representing the REST resource this\nobject represents. Servers may infer this from the endpoint the client\nsubmits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
	//									"type":        "string",
	//								},
	//								"metadata": map[string]interface{}{
	//									"type": "object",
	//								},
	//								"spec": map[string]interface{}{
	//									"description": "SeataServerSpec defines the desired state of SeataServer",
	//									"properties": map[string]interface{}{
	//										"containerName": "seata-server",
	//										"type":          "string",
	//									},
	//									"env": map[string]interface{}{
	//										"additionalProperties": map[string]interface{}{
	//											"type": "string",
	//										},
	//										"type": "object",
	//									},
	//									"image": map[string]interface{}{
	//										"type": "string",
	//									},
	//									"ports": map[string]interface{}{
	//										"properties": map[string]interface{}{
	//											"consolePort:": map[string]interface{}{
	//												"default": "7091",
	//												"format":  "int32",
	//												"type":    "integer",
	//											},
	//											"raftPort:": map[string]interface{}{
	//												"default": "9091",
	//												"format":  "int32",
	//												"type":    "integer",
	//											},
	//											"servicePort:": map[string]interface{}{
	//												"default": "89091",
	//												"format":  "int32",
	//												"type":    "integer",
	//											},
	//											"type": "string",
	//										},
	//									},
	//								},
	//								"served":  true,
	//								"storage": true,
	//							},
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}
	// 获取 GVR（Group Version Resource）
	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}

	// 创建 CRD
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
