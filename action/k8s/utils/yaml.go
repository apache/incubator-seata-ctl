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

import (
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

// ConvertAndSaveYamlToJSON takes a YAML file path and a target file name, converts it to JSON,
// and does not convert if the JSON file already exists.
func ConvertAndSaveYamlToJSON(targetName string) (string, error) {
	// Get the original file directory, filename, and extension
	dir, file := filepath.Split(targetName)
	fileName := file[:len(file)-len(filepath.Ext(file))] // Remove the original file extension
	if targetName == "" {
		targetName = fileName // If no target name is provided, use the original YAML file name
	}
	newFilePath := filepath.Join(dir, fileName+".json") // Generate the new JSON file path

	// Check if the JSON file already exists
	if _, err := os.Stat(newFilePath); err == nil {
		return newFilePath, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to check JSON file existence: %v", err)
	}

	// Read the YAML file
	yamlData, err := os.ReadFile(targetName)
	if err != nil {
		return "", fmt.Errorf("failed to read YAML file: %v", err)
	}

	// Convert the YAML to JSON
	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return "", fmt.Errorf("failed to convert YAML to JSON: %v", err)
	}

	// Write the converted JSON data to the new file
	err = os.WriteFile(newFilePath, jsonData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write JSON file: %v", err)
	}
	return newFilePath, nil
}
