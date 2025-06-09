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

package logadapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/seata/seata-ctl/tool"
)

// QueryLogs sends a request to the /query endpoint and retrieves logs based on the provided filter.
// - filter: A map that holds log filtering parameters such as log level.
// - currency: A struct containing information like the source application and API address.
// - number: The number of logs to fetch (limit).
func (l *Local) QueryLogs(filter map[string]interface{}, currency *Currency, number int) error {
	logLevel, ok := filter["logLevel"].(string)
	if !ok {
		return fmt.Errorf("logLevel is missing or invalid in the filter")
	}

	// Build the query URL using the filter and currency information
	url := fmt.Sprintf("%s%s?application_id=%s&log_level=%s&limit=%d", currency.Address, LocalQueryPath, currency.Source, logLevel, number)

	// Send a GET request to the /query endpoint
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make request to %s: %v", url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			tool.Logger.Errorf("failed to close response body: %v", err)
		}
	}(resp.Body)

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Check if the response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	// Parse the JSON response body into QueryResponse struct
	var queryResponse QueryResponse
	err = json.Unmarshal(body, &queryResponse)
	if err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Print the logs in a structured format
	for _, value := range queryResponse.Logs {
		res := fmt.Sprintf("[%s]: %s\n", value.Timestamp, value.LogMessage)
		fmt.Println(value.LogMessage)
		if strings.Contains(value.LogMessage, "INFO") {
			tool.Logger.Info(fmt.Sprintf("%v", res))
		}
		if strings.Contains(value.LogMessage, "ERROR") {
			tool.Logger.Error(fmt.Sprintf("%v", res))
		}
		if strings.Contains(value.LogMessage, "WARN") {
			tool.Logger.Warn(fmt.Sprintf("%v", res))
		}
	}
	return nil
}
