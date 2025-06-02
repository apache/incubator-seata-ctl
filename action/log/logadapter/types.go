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

import "time"

type Elasticsearch struct{}

type Loki struct{}

type Local struct{}

type Currency struct {
	Address  string `json:"address"`
	Source   string `json:"source"`
	Username string `json:"username"`
	Password string `json:"password"`
	Index    string `json:"index"`
}

type SeataLog struct {
	Timestamp  time.Time // Timestamp
	LogLevel   string    // Log level, e.g., INFO, ERROR, etc.
	Module     string    // Log module, e.g., RMHandler, UndoLogManager, etc.
	XID        string    // Global transaction ID
	BranchID   string    // Branch transaction ID
	ResourceID string    // Resource ID
	Message    string    // Log content
}

// LokiQueryResult represents the structure of the response from Loki
type LokiQueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// LogData represents a single log entry
type LogData struct {
	Timestamp  string `json:"timestamp"`
	LogLevel   string `json:"log_level"`
	LogMessage string `json:"log_message"`
}

// QueryResponse holds the response structure from the /query API
type QueryResponse struct {
	ApplicationID string    `json:"application_id"`
	LogLevel      string    `json:"log_level"`
	Logs          []LogData `json:"logs"`
}
