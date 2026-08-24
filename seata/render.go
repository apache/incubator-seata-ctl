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

package seata

import (
	"encoding/json"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

const (
	OutputTable = "table"
	OutputJSON  = "json"
	OutputYAML  = "yaml"
)

func NormalizeOutput(output string) (string, error) {
	output = strings.ToLower(output)
	switch output {
	case OutputTable, OutputJSON, OutputYAML:
		return output, nil
	default:
		return "", fmt.Errorf("unsupported output format %q", output)
	}
}

func FormatStructuredOutput(value interface{}, output string) (string, error) {
	switch output {
	case OutputJSON:
		data, err := json.MarshalIndent(value, "", "  ")
		return string(data), err
	case OutputYAML:
		data, err := yaml.Marshal(value)
		return string(data), err
	default:
		return "", fmt.Errorf("unsupported structured output format %q", output)
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func intValue(value *int) interface{} {
	if value == nil {
		return ""
	}
	return *value
}

func int64Value(value *int64) interface{} {
	if value == nil {
		return ""
	}
	return *value
}
