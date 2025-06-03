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

	"github.com/seata/seata-ctl/action/k8s/utils"
	"github.com/seata/seata-ctl/tool"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var UnDeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "undeploy seata in k8s",
	Run: func(_ *cobra.Command, _ []string) {
		err := undeploy()
		if err != nil {
			tool.Logger.Errorf("undeploy error: %v", err)
		}
	},
}

func init() {
	UnDeployCmd.PersistentFlags().StringVar(&Name, "name", DefaultCRName, "Seataserver name")
	UnDeployCmd.PersistentFlags().StringVar(&Namespace, "namespace", DefaultNamespace, "Namespace name")
}

func undeploy() error {
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

	err = client.Resource(gvr).Namespace(namespace).Delete(context.TODO(), Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	tool.Logger.Infof("CR delete success，name: %s\n", Name)

	return nil
}
