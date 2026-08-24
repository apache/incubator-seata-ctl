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
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestQueryGlobalLocksBuildsRequest(t *testing.T) {
	timeStart := int64(1710000000000)
	timeEnd := int64(1710003600000)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != GlobalLockQueryURL {
			t.Fatalf("path = %s, want %s", r.URL.Path, GlobalLockQueryURL)
		}
		if got := r.Header.Get("authorization"); got != "test-token" {
			t.Fatalf("authorization = %q, want test-token", got)
		}
		wantQuery := map[string]string{
			"xid":           "xid-1",
			"tableName":     "account",
			"transactionId": "1001",
			"branchId":      "2001",
			"pk":            "1",
			"resourceId":    "jdbc:mysql://127.0.0.1:3306/seata",
			"pageNum":       "2",
			"pageSize":      "10",
			"timeStart":     "1710000000000",
			"timeEnd":       "1710003600000",
		}
		for key, want := range wantQuery {
			if got := r.URL.Query().Get(key); got != want {
				t.Fatalf("query[%s] = %q, want %q", key, got, want)
			}
		}
		fmt.Fprint(w, `{"code":"200","message":"success","success":true,"pageSize":10,"pageNum":2,"total":1,"pages":1,"data":[{"xid":"xid-1","transactionId":"1001","branchId":"2001","resourceId":"jdbc:mysql://127.0.0.1:3306/seata","tableName":"account","pk":"1","rowKey":"1","vgroup":"default_tx_group"}]}`)
	}))
	defer server.Close()
	defer setTestAuth(t, server.URL, "test-token")()

	response, err := queryGlobalLocks(&ConsoleClient{HTTPClient: server.Client()}, GlobalLockQuery{
		XID:           "xid-1",
		TableName:     "account",
		TransactionID: "1001",
		BranchID:      "2001",
		PK:            "1",
		ResourceID:    "jdbc:mysql://127.0.0.1:3306/seata",
		PageNum:       2,
		PageSize:      10,
		TimeStart:     &timeStart,
		TimeEnd:       &timeEnd,
	})
	if err != nil {
		t.Fatalf("queryGlobalLocks returned error: %v", err)
	}
	if len(response.Data) != 1 || stringValue(response.Data[0].TableName) != "account" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestCheckGlobalLockBuildsRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != GlobalLockCheckURL {
			t.Fatalf("path = %s, want %s", r.URL.Path, GlobalLockCheckURL)
		}
		if got := r.URL.Query().Get("xid"); got != "xid-1" {
			t.Fatalf("xid = %q, want xid-1", got)
		}
		if got := r.URL.Query().Get("branchId"); got != "2001" {
			t.Fatalf("branchId = %q, want 2001", got)
		}
		fmt.Fprint(w, `{"code":"200","message":"success","success":true,"data":true}`)
	}))
	defer server.Close()
	defer setTestAuth(t, server.URL, "test-token")()

	response, err := checkGlobalLock(&ConsoleClient{HTTPClient: server.Client()}, "xid-1", "2001")
	if err != nil {
		t.Fatalf("checkGlobalLock returned error: %v", err)
	}
	if response.Data == nil || !*response.Data {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestQueryGlobalLocksErrors(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "server code", body: `{"code":"500","message":"boom"}`, want: "boom"},
		{name: "invalid json", body: `{`, want: "decode /api/v1/console/globalLock/query response"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, tt.body)
			}))
			defer server.Close()
			defer setTestAuth(t, server.URL, "test-token")()

			_, err := queryGlobalLocks(&ConsoleClient{HTTPClient: server.Client()}, GlobalLockQuery{PageNum: 1, PageSize: 20})
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want containing %q", err, tt.want)
			}
		})
	}
}

func TestFormatGlobalLockPageAndCheck(t *testing.T) {
	success := true
	response := &GlobalLockPageResult{
		Code:    CodeOK,
		Message: "success",
		Success: &success,
		PageNum:  intPtr(1),
		PageSize: intPtr(20),
		Total:    intPtr(1),
		Pages:    intPtr(1),
		Data: []GlobalLock{{
			XID:           stringPtr("xid-1"),
			TransactionID: stringPtr("1001"),
			BranchID:      stringPtr("2001"),
			ResourceID:    stringPtr("jdbc:mysql://127.0.0.1:3306/seata"),
			TableName:     stringPtr("account"),
			PK:            stringPtr("1"),
			RowKey:        stringPtr("1"),
			Vgroup:        stringPtr("default_tx_group"),
		}},
	}

	tableOutput, err := FormatGlobalLockPage(response, OutputTable)
	if err != nil {
		t.Fatalf("table output error: %v", err)
	}
	for _, want := range []string{"xid", "transaction_id", "branch_id", "resource_id", "table_name", "pk", "row_key", "vgroup"} {
		if !strings.Contains(tableOutput, want) {
			t.Fatalf("table output %q does not contain %q", tableOutput, want)
		}
	}

	jsonOutput, err := FormatGlobalLockPage(response, OutputJSON)
	if err != nil {
		t.Fatalf("json output error: %v", err)
	}
	var jsonResult GlobalLockPageResult
	if err = json.Unmarshal([]byte(jsonOutput), &jsonResult); err != nil {
		t.Fatalf("unmarshal json output: %v", err)
	}
	if !reflect.DeepEqual(&jsonResult, response) {
		t.Fatalf("unexpected json result: %+v", jsonResult)
	}

	yamlOutput, err := FormatGlobalLockPage(response, OutputYAML)
	if err != nil {
		t.Fatalf("yaml output error: %v", err)
	}
	var yamlResult GlobalLockPageResult
	if err = yaml.Unmarshal([]byte(yamlOutput), &yamlResult); err != nil {
		t.Fatalf("unmarshal yaml output: %v", err)
	}
	if len(yamlResult.Data) != 1 || stringValue(yamlResult.Data[0].TableName) != "account" {
		t.Fatalf("unexpected yaml result: %+v", yamlResult)
	}

	checkOutput, err := FormatGlobalLockCheck(&GlobalLockCheckResult{Code: CodeOK, Message: "locked", Data: &success}, OutputTable)
	if err != nil {
		t.Fatalf("check table output error: %v", err)
	}
	for _, want := range []string{"locked", "message"} {
		if !strings.Contains(checkOutput, want) {
			t.Fatalf("check output %q does not contain %q", checkOutput, want)
		}
	}

	_, err = FormatGlobalLockPage(response, "xml")
	if err == nil || !strings.Contains(err.Error(), "unsupported output format") {
		t.Fatalf("error = %v, want unsupported output format", err)
	}
}

func TestCheckGlobalLockValidation(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("should not be called")
	})}
	_, err := checkGlobalLock(&ConsoleClient{HTTPClient: client}, "", "2001")
	if err == nil || !strings.Contains(err.Error(), "xid is required") {
		t.Fatalf("error = %v, want xid is required", err)
	}
	_, err = checkGlobalLock(&ConsoleClient{HTTPClient: client}, "xid-1", "")
	if err == nil || !strings.Contains(err.Error(), "branch-id is required") {
		t.Fatalf("error = %v, want branch-id is required", err)
	}
}
