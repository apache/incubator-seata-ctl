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
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "show seata status in k8s",
	Run: func(_ *cobra.Command, _ []string) {
		err := status()
		if err != nil {
			tool.Logger.Errorf("get k8s status error: %v", err)
		}
	},
}

func init() {
	StatusCmd.PersistentFlags().StringVar(&Name, "name", "list", "Seataserver name")
	StatusCmd.PersistentFlags().StringVar(&Namespace, "namespace", DefaultNamespace, "Namespace name")
}

func status() error {
	statuses, err := getPodsStatusByLabel(Namespace, Name)
	if err != nil {
		return err
	}
	// Print formatted Pod status information
	for _, status := range statuses {
		tool.Logger.Infof("status: %s", status)
	}
	return nil
}

func getPodsStatusByLabel(namespace, labelSelector string) ([]string, error) {
	client, err := utils.GetClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	// Retrieve all Pods with the specified label
	// Use LabelSelector to filter Pods with the specified cr_name label
	labelSelector = Label + "=" + labelSelector
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	// Check if any pods were found
	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no matching pods found") // Alternatively, return a specific error message as needed
	}

	// Iterate over all Pods and get their status
	var statuses []string

	// Build formatted status string for output
	statuses = append(statuses, fmt.Sprintf("%-25s %-10s", "POD NAME", "STATUS")) // Header
	statuses = append(statuses, "-------------------------------------------")

	for _, pod := range pods.Items {
		statuses = append(statuses, fmt.Sprintf("%-25s %-10s", pod.Name, pod.Status.Phase))
	}

	return statuses, nil
}
