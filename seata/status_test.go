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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestQueryStatusBuildsRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != HealthCheckURL {
			t.Fatalf("path = %s, want %s", r.URL.Path, HealthCheckURL)
		}
		if got := r.Header.Get("authorization"); got != "test-token" {
			t.Fatalf("authorization = %q, want test-token", got)
		}
		fmt.Fprint(w, `{"code":"200","message":"success","success":true,"data":[{"type":"nacos","address":"127.0.0.1:7091","status":"ok"}]}`)
	}))
	defer server.Close()
	defer setTestAuth(t, server.URL, "test-token")()

	response, err := queryStatus(&ConsoleClient{HTTPClient: server.Client()})
	if err != nil {
		t.Fatalf("queryStatus returned error: %v", err)
	}
	if len(response.Data) != 1 || response.Data[0].Type != "nacos" || response.Data[0].Status != "ok" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestFormatNodeStatusResponse(t *testing.T) {
	response := &NodeStatusResponse{
		BaseResponse: BaseResponse{Code: CodeOK, Message: "success", Success: true},
		Data: []NodeStatus{{Type: "nacos", Address: "127.0.0.1:7091", Status: "ok"}},
	}

	tableOutput, err := FormatNodeStatusResponse(response, OutputTable)
	if err != nil {
		t.Fatalf("table output error: %v", err)
	}
	for _, want := range []string{"type", "address", "status", "nacos", "127.0.0.1:7091", "ok"} {
		if !strings.Contains(tableOutput, want) {
			t.Fatalf("table output %q does not contain %q", tableOutput, want)
		}
	}

	jsonOutput, err := FormatNodeStatusResponse(response, OutputJSON)
	if err != nil {
		t.Fatalf("json output error: %v", err)
	}
	var jsonResult NodeStatusResponse
	if err = json.Unmarshal([]byte(jsonOutput), &jsonResult); err != nil {
		t.Fatalf("unmarshal json output: %v", err)
	}
	if len(jsonResult.Data) != 1 || jsonResult.Data[0].Address != "127.0.0.1:7091" {
		t.Fatalf("unexpected json result: %+v", jsonResult)
	}

	yamlOutput, err := FormatNodeStatusResponse(response, OutputYAML)
	if err != nil {
		t.Fatalf("yaml output error: %v", err)
	}
	var yamlResult NodeStatusResponse
	if err = yaml.Unmarshal([]byte(yamlOutput), &yamlResult); err != nil {
		t.Fatalf("unmarshal yaml output: %v", err)
	}
	if len(yamlResult.Data) != 1 || yamlResult.Data[0].Type != "nacos" {
		t.Fatalf("unexpected yaml result: %+v", yamlResult)
	}

	_, err = FormatNodeStatusResponse(response, "xml")
	if err == nil || !strings.Contains(err.Error(), "unsupported output format") {
		t.Fatalf("error = %v, want unsupported output format", err)
	}
}
