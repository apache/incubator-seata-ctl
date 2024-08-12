package k8s

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/seata/seata-ctl/action/common"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "show seata status in k8s",
	Run: func(cmd *cobra.Command, args []string) {
		err := status()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	StatusCmd.PersistentFlags().StringVar(&Name, "name", "list", "Seataserver name")
	StatusCmd.PersistentFlags().StringVar(&Namespace, "namespace", "default", "Namespace name")
}

func status() error {
	//获取动态kubeclient
	client, err := common.GetDynamicClient()
	if err != nil {
		return err
	}
	// 获取命名空间
	namespace := Namespace

	// 定义 Custom Resource 的 GroupVersionResource
	gvr := schema.GroupVersionResource{
		Group:    "operator.seata.apache.org",
		Version:  "v1alpha1",
		Resource: "seataservers",
	}

	var jsonData []byte
	if Name == "list" {
		fmt.Printf("1")
		var seataServerList *unstructured.UnstructuredList
		seataServerList, err = client.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		jsonData, err = seataServerList.MarshalJSON()
		if err != nil {
			log.Fatalf("Error marshalling Unstructured object to JSON: %v", err)
		}
	} else {
		var seataServer *unstructured.Unstructured
		seataServer, err = client.Resource(gvr).Namespace(namespace).Get(context.TODO(), Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		jsonData, err = seataServer.MarshalJSON()
		if err != nil {
			log.Fatalf("Error marshalling Unstructured object to JSON: %v", err)
		}
	}
	//更美观的格式化输出
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, jsonData, "", "")
	if err != nil {
		log.Fatalf("Error indenting JSON: %v", err)
	}
	fmt.Printf("%s\n", prettyJSON.Bytes())
	return nil
}
