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
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestQueryGlobalSessionsBuildsRequest(t *testing.T) {
	status := 1
	withBranch := true
	timeStart := int64(1710000000000)
	timeEnd := int64(1710003600000)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != GlobalSessionQueryURL {
			t.Fatalf("path = %s, want %s", r.URL.Path, GlobalSessionQueryURL)
		}
		if got := r.Header.Get("authorization"); got != "test-token" {
			t.Fatalf("authorization = %q, want test-token", got)
		}
		wantQuery := map[string]string{
			"xid":             "xid-1",
			"applicationId":   "account-service",
			"status":          "1",
			"transactionName": "create-order",
			"vgroup":          "default_tx_group",
			"withBranch":      "true",
			"pageNum":         "2",
			"pageSize":        "10",
			"timeStart":       "1710000000000",
			"timeEnd":         "1710003600000",
		}
		for key, want := range wantQuery {
			if got := r.URL.Query().Get(key); got != want {
				t.Fatalf("query[%s] = %q, want %q", key, got, want)
			}
		}
		fmt.Fprint(w, `{
			"code":"200",
			"message":"success",
			"success":true,
			"pageSize":10,
			"pageNum":2,
			"total":1,
			"pages":1,
			"data":[{
				"xid":"xid-1",
				"transactionId":"1001",
				"status":1,
				"applicationId":"account-service",
				"transactionServiceGroup":"default_tx_group",
				"transactionName":"create-order",
				"timeout":30000,
				"beginTime":1710000000000
			}]
		}`)
	}))
	defer server.Close()
	defer setTestAuth(t, server.URL, "test-token")()

	response, err := queryGlobalSessions(&ConsoleClient{HTTPClient: server.Client()}, GlobalSessionQuery{
		XID:             "xid-1",
		ApplicationID:   "account-service",
		Status:          &status,
		TransactionName: "create-order",
		Vgroup:          "default_tx_group",
		WithBranch:      &withBranch,
		PageNum:         2,
		PageSize:        10,
		TimeStart:       &timeStart,
		TimeEnd:         &timeEnd,
	})
	if err != nil {
		t.Fatalf("queryGlobalSessions returned error: %v", err)
	}
	if len(response.Data) != 1 {
		t.Fatalf("len(response.Data) = %d, want 1", len(response.Data))
	}
	session := response.Data[0]
	if stringValue(session.XID) != "xid-1" || stringValue(session.TransactionID) != "1001" || stringValue(session.TransactionServiceGroup) != "default_tx_group" {
		t.Fatalf("unexpected session: %+v", session)
	}
}

func TestQueryGlobalSessionsOmitsUnsetFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, key := range []string{"xid", "applicationId", "status", "transactionName", "vgroup", "withBranch", "timeStart", "timeEnd"} {
			if _, ok := r.URL.Query()[key]; ok {
				t.Fatalf("query[%s] should be omitted", key)
			}
		}
		if got := r.URL.Query().Get("pageNum"); got != "1" {
			t.Fatalf("pageNum = %q, want 1", got)
		}
		if got := r.URL.Query().Get("pageSize"); got != "20" {
			t.Fatalf("pageSize = %q, want 20", got)
		}
		fmt.Fprint(w, `{"code":"200","message":"success","data":[]}`)
	}))
	defer server.Close()
	defer setTestAuth(t, server.URL, "test-token")()

	_, err := queryGlobalSessions(&ConsoleClient{HTTPClient: server.Client()}, GlobalSessionQuery{PageNum: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("queryGlobalSessions returned error: %v", err)
	}
}

func TestQueryGlobalSessionsErrors(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "server code", body: `{"code":"500","message":"boom"}`, want: "boom"},
		{name: "invalid json", body: `{`, want: "decode /api/v1/console/globalSession/query response"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, tt.body)
			}))
			defer server.Close()
			defer setTestAuth(t, server.URL, "test-token")()

			_, err := queryGlobalSessions(&ConsoleClient{HTTPClient: server.Client()}, GlobalSessionQuery{PageNum: 1, PageSize: 20})
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want containing %q", err, tt.want)
			}
		})
	}
}

func TestQueryGlobalSessionsHTTPAndNetworkErrors(t *testing.T) {
	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "no route", http.StatusNotFound)
		}))
		defer server.Close()
		defer setTestAuth(t, server.URL, "test-token")()

		_, err := queryGlobalSessions(&ConsoleClient{HTTPClient: server.Client()}, GlobalSessionQuery{PageNum: 1, PageSize: 20})
		if err == nil || !strings.Contains(err.Error(), "http status 404") {
			t.Fatalf("error = %v, want http status 404", err)
		}
	})

	t.Run("network", func(t *testing.T) {
		oldAuth := auth
		auth = Auth{ServerIP: "127.0.0.1", ServerPort: 7091, token: "test-token"}
		defer func() { auth = oldAuth }()

		client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("dial failed")
		})}
		_, err := queryGlobalSessions(&ConsoleClient{HTTPClient: client}, GlobalSessionQuery{PageNum: 1, PageSize: 20})
		if err == nil || !strings.Contains(err.Error(), "dial failed") {
			t.Fatalf("error = %v, want dial failed", err)
		}
	})
}

func TestFormatGlobalSessionPage(t *testing.T) {
	success := true
	response := &GlobalSessionPageResult{
		Code:     CodeOK,
		Message:  "success",
		Success:  &success,
		PageNum:  intPtr(1),
		PageSize: intPtr(20),
		Total:    intPtr(1),
		Pages:    intPtr(1),
		Data: []GlobalSession{{
			XID:                     stringPtr("xid-1"),
			TransactionID:           stringPtr("1001"),
			Status:                  intPtr(1),
			ApplicationID:           stringPtr("account-service"),
			TransactionServiceGroup: stringPtr("default_tx_group"),
			TransactionName:         stringPtr("create-order"),
			BeginTime:               int64Ptr(1710000000000),
			Timeout:                 int64Ptr(30000),
		}},
	}

	tableOutput, err := FormatGlobalSessionPage(response, OutputTable)
	if err != nil {
		t.Fatalf("table output error: %v", err)
	}
	for _, want := range []string{
		"| xid   | transaction_id | status | application_id  | vgroup           | transaction_name |    begin_time | timeout |",
		"| xid-1 | 1001           |      1 | account-service | default_tx_group | create-order     | 1710000000000 |   30000 |",
	} {
		if !strings.Contains(tableOutput, want) {
			t.Fatalf("table output %q does not contain %q", tableOutput, want)
		}
	}

	jsonOutput, err := FormatGlobalSessionPage(response, OutputJSON)
	if err != nil {
		t.Fatalf("json output error: %v", err)
	}
	var jsonResult GlobalSessionPageResult
	if err = json.Unmarshal([]byte(jsonOutput), &jsonResult); err != nil {
		t.Fatalf("unmarshal json output: %v", err)
	}
	if !reflect.DeepEqual(&jsonResult, response) {
		t.Fatalf("unexpected json result: %+v", jsonResult)
	}
	for _, want := range []string{`"applicationData": null`, `"branchSessionVOs": null`} {
		if !strings.Contains(jsonOutput, want) {
			t.Fatalf("json output %q does not contain %q", jsonOutput, want)
		}
	}

	yamlOutput, err := FormatGlobalSessionPage(response, OutputYAML)
	if err != nil {
		t.Fatalf("yaml output error: %v", err)
	}
	var yamlResult GlobalSessionPageResult
	if err = yaml.Unmarshal([]byte(yamlOutput), &yamlResult); err != nil {
		t.Fatalf("unmarshal yaml output: %v", err)
	}
	if yamlResult.PageNum == nil || *yamlResult.PageNum != 1 || len(yamlResult.Data) != 1 || stringValue(yamlResult.Data[0].XID) != "xid-1" || yamlResult.Data[0].ApplicationData != nil {
		t.Fatalf("unexpected yaml result: %+v", yamlResult)
	}
	for _, want := range []string{"applicationData: null", "branchSessionVOs: null"} {
		if !strings.Contains(yamlOutput, want) {
			t.Fatalf("yaml output %q does not contain %q", yamlOutput, want)
		}
	}

	_, err = FormatGlobalSessionPage(response, "xml")
	if err == nil || !strings.Contains(err.Error(), "unsupported output format") {
		t.Fatalf("error = %v, want unsupported output format", err)
	}
}

func setTestAuth(t *testing.T, serverURL string, token string) func() {
	t.Helper()
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		t.Fatalf("parse server url: %v", err)
	}
	host, portStr, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		t.Fatalf("split host port: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parse port: %v", err)
	}
	oldAuth := auth
	auth = Auth{ServerIP: host, ServerPort: port, token: token}
	return func() {
		auth = oldAuth
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func stringPtr(value string) *string {
	return &value
}

func intPtr(value int) *int {
	return &value
}

func int64Ptr(value int64) *int64 {
	return &value
}
