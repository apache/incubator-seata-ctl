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
	"github.com/seata/seata-ctl/tool"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// QueryLogs queries logs from Loki based on the filter and settings provided
func (l *Loki) QueryLogs(filter map[string]interface{}, currency *Currency, number int) error {

	// Prepare the query URL with time range
	params := url.Values{}
	params.Set("query", filter["query"].(string))
	params.Set("limit", strconv.Itoa(number))

	// Set start time if provided
	if value, ok := filter["start"]; ok {
		res, err := parseToTimestamp(value.(string))
		if err != nil {
			return err
		}
		params.Set("start", fmt.Sprintf("%d", res))
	}
	// Set end time if provided
	if value, ok := filter["end"]; ok {
		res, err := parseToTimestamp(value.(string))
		if err != nil {
			return err
		}
		params.Set("end", fmt.Sprintf("%d", res))
	}
	queryURL := currency.Address + LokiAddressPath + params.Encode()

	// Send GET request to Loki
	resp, err := http.Get(queryURL)
	if err != nil {
		return fmt.Errorf("error sending request to Loki: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			tool.Logger.Errorf("Failed to close response body: %v", err)
		}
	}(resp.Body)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response from Loki: %v", err)
	}

	// Parse the JSON response
	var result LokiQueryResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("error parsing Loki response: %v", err)
	}

	// Process and print logs
	if result.Status == "success" {
		if len(result.Data.Result) == 0 {
			return fmt.Errorf("loki query returned no results")
		}
		for _, stream := range result.Data.Result {
			for _, entry := range stream.Values {
				// Extract timestamp and log message
				value := entry[1]
				// Print the readable timestamp and log message
				if strings.Contains(value, "INFO") {
					tool.Logger.Info(fmt.Sprintf("%v", value))
				}
				if strings.Contains(value, "ERROR") {
					tool.Logger.Error(fmt.Sprintf("%v", value))
				}
				if strings.Contains(value, "WARN") {
					tool.Logger.Warn(fmt.Sprintf("%v", value))
				}
			}
		}
	} else {
		fmt.Printf("Query failed, status: %s\n", result.Status)
	}
	return nil
}

// parseToTimestamp parses a time string to a Unix nanosecond timestamp
func parseToTimestamp(timeStr string) (int64, error) {
	// Load the specified time zone (UTC in this case)
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return 0, fmt.Errorf("failed to load timezone: %v", err)
	}

	// Parse the time string using the specified layout and timezone
	t, err := time.ParseInLocation(TimeLayout, timeStr, loc)
	if err != nil {
		return 0, fmt.Errorf("failed to parse time: %v", err)
	}

	// Convert to Unix nanosecond timestamp
	timestamp := t.UnixNano()
	return timestamp, nil
}
