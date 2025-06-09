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

package utils

import "fmt"

func CreateRequest(createCrdPath string, filePath string) (string, error) {
	kubeConfigPath, err := GetKubeConfigPath()
	if err != nil {
		return "", fmt.Errorf("failed to get kubeconfig path: %v", err)
	}
	kubeConfig, err := LoadKubeConfig(kubeConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to load kubeconfig file: %v", err)
	}
	contextInfo, err := GetContextInfo(kubeConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get context info: %v", err)
	}

	filePath, err = ConvertAndSaveYamlToJSON(filePath)

	if err != nil {
		return "", fmt.Errorf("failed to save yaml: %v", err)
	}

	res, err := sendPostRequest(contextInfo, createCrdPath, filePath)
	if err != nil {
		return res, fmt.Errorf("failed to send post request: %v", err)
	}
	return res, nil
}
