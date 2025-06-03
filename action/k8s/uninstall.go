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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var UnInstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall seata in k8s",
	Run: func(_ *cobra.Command, _ []string) {
		err := UninstallCRD()
		if err != nil {
			tool.Logger.Errorf("uninstall CRD err:%v", err)
		}
		err = UnDeploymentController()
		if err != nil {
			tool.Logger.Errorf("uninstall Deployment err:%v", err)
		}
	},
}

func init() {
	UnInstallCmd.PersistentFlags().StringVar(&Namespace, "namespace", DefaultNamespace, "Namespace name")
	UnInstallCmd.PersistentFlags().StringVar(&DeployName, "name", DefaultDeployName, "Deployment name")
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
			tool.Logger.Errorf("CRD %s does not exist, no action taken.\n", CRDname)
		} else {
			// For other errors, log the error and exit the program
			tool.Logger.Errorf("delete CRD failed: %v", err)
		}
	} else {
		// Successfully deleted the resource
		tool.Logger.Infof("delete CRD %s successfully.\n", CRDname)
	}

	return nil
}

func UnDeploymentController() error {
	client, err := utils.GetClient()
	if err != nil {
		return err
	}

	// Assume client has already been defined
	err = client.AppsV1().Deployments(Namespace).Delete(context.TODO(), DeployName, metav1.DeleteOptions{})
	if err != nil {
		// Check if the error is a "not found" error
		if errors.IsNotFound(err) {
			// The deployment does not exist, output a message instead of returning an error
			tool.Logger.Errorf("Deployment '%s' does not exist in namespace '%s', no action taken.\n", DeployName, Namespace)
		} else {
			// For other errors, log the error and exit the program
			tool.Logger.Errorf("Error deleting deployment: %s", err.Error())
		}
	} else {
		// Successfully deleted the deployment
		tool.Logger.Infof("deleted Controller %s successfully ", DeployName)
	}
	return nil
}
