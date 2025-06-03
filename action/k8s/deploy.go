/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package k8s

import (
	"context"
	"fmt"

	"github.com/seata/seata-ctl/action/k8s/utils"
	"github.com/seata/seata-ctl/tool"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy seata in k8s",
	Run: func(_ *cobra.Command, _ []string) {
		err := deploy()
		if err != nil {
			tool.Logger.Errorf("deploy err:%v", err)
		}
	},
}

func init() {
	DeployCmd.PersistentFlags().StringVar(&Name, "name", DefaultCRName, "Seataserver name")
	DeployCmd.PersistentFlags().Int32Var(&Replicas, "replicas", 1, "Replicas number")
	DeployCmd.PersistentFlags().StringVar(&Namespace, "namespace", DefaultNamespace, "Namespace name")
	DeployCmd.PersistentFlags().StringVar(&Image, "image", DefaultServerImage, "Seata server image")
}

func deploy() error {
	client, err := utils.GetDynamicClient()
	if err != nil {
		return err
	}
	namespace := Namespace
	gvr := schema.GroupVersionResource{
		Group:    "operator.seata.apache.org",
		Version:  "v1alpha1",
		Resource: "seataservers",
	}

	var seataServer *unstructured.Unstructured
	seataServer, err = client.Resource(gvr).Namespace(namespace).Get(context.TODO(), Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if seataServer != nil {
		return fmt.Errorf("seata server already exist! name:" + Name)
	}
	seataServer = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "operator.seata.apache.org/v1alpha1",
			"kind":       "SeataServer",
			"metadata": map[string]interface{}{
				"name":      Name,
				"namespace": Namespace,
			},
			"spec": map[string]interface{}{
				"serviceName": ServiceName,
				"replicas":    Replicas,
				"image":       Image,
				"store": map[string]interface{}{
					"resources": map[string]interface{}{
						"requests": map[string]interface{}{
							"storage": RequestStorage,
						},
						"limits": map[string]interface{}{
							"storage": LimitStorage,
						},
					},
				},
			},
		},
	}
	_, err = client.Resource(gvr).Namespace(namespace).Create(context.TODO(), seataServer, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	tool.Logger.Infof("CR install success，name: %s\n", Name)
	return nil
}
