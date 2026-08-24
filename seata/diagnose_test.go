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
	"reflect"
	"strings"
	"testing"
	"time"

	yaml "gopkg.in/yaml.v3"
)

func TestRunDiagnostics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("authorization"); got != "test-token" {
			t.Fatalf("authorization = %q, want test-token", got)
		}
		switch r.URL.Path {
		case HealthCheckURL:
			fmt.Fprint(w, `{"code":"200","message":"success","success":true,"data":[{"type":"nacos","address":"127.0.0.1:7091","status":"ok"}]}`)
		case GlobalSessionQueryURL:
			if got := r.URL.Query().Get("pageNum"); got != "1" {
				t.Fatalf("pageNum = %q, want 1", got)
			}
			if got := r.URL.Query().Get("pageSize"); got != "1" {
				t.Fatalf("pageSize = %q, want 1", got)
			}
			fmt.Fprint(w, `{"code":"200","message":"success","success":true,"pageSize":1,"pageNum":1,"total":1,"pages":1,"data":[{"xid":"xid-1","transactionId":"1001","status":1,"applicationId":"order-service","transactionServiceGroup":"default_tx_group","transactionName":"createOrder","timeout":60000,"beginTime":1710000000000}]}`)
		case GlobalLockQueryURL:
			if got := r.URL.Query().Get("pageNum"); got != "1" {
				t.Fatalf("lock pageNum = %q, want 1", got)
			}
			if got := r.URL.Query().Get("pageSize"); got != "1" {
				t.Fatalf("lock pageSize = %q, want 1", got)
			}
			fmt.Fprint(w, `{"code":"200","message":"success","success":true,"pageSize":1,"pageNum":1,"total":1,"pages":1,"data":[{"xid":"xid-1","transactionId":"1001","branchId":"2001","resourceId":"jdbc:mysql://127.0.0.1:3306/seata","tableName":"account","pk":"1","rowKey":"1","vgroup":"default_tx_group"}]}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	defer setTestAuth(t, server.URL, "test-token")()

	report, err := RunDiagnostics(DiagnoseOptions{
		TCPTimeout:   time.Second,
		PageNum:      1,
		PageSize:     1,
		LockPageNum:  1,
		LockPageSize: 1,
	})
	if err != nil {
		t.Fatalf("RunDiagnostics returned error: %v", err)
	}
	if report.Success == nil || !*report.Success {
		t.Fatalf("unexpected success flag: %+v", report.Success)
	}
	wantStages := []string{"config", "tcp connectivity", "login token", "status", "transaction query", "global lock query", "db/schema"}
	if len(report.Data) != len(wantStages) {
		t.Fatalf("len(report.Data) = %d, want %d", len(report.Data), len(wantStages))
	}
	for i, want := range wantStages {
		if report.Data[i].Stage != want {
			t.Fatalf("stage[%d] = %q, want %q", i, report.Data[i].Stage, want)
		}
	}
	for _, stage := range report.Data[:6] {
		if stage.Status != DiagnoseStatusPass {
			t.Fatalf("stage %+v not pass", stage)
		}
	}
	if report.Data[6].Status != DiagnoseStatusSkip {
		t.Fatalf("db stage = %+v, want skip", report.Data[6])
	}

	tableOutput, err := FormatDiagnoseReport(report, OutputTable)
	if err != nil {
		t.Fatalf("table output error: %v", err)
	}
	for _, want := range []string{"config", "tcp connectivity", "transaction query", "global lock query", "db/schema"} {
		if !strings.Contains(tableOutput, want) {
			t.Fatalf("table output %q does not contain %q", tableOutput, want)
		}
	}

	jsonOutput, err := FormatDiagnoseReport(report, OutputJSON)
	if err != nil {
		t.Fatalf("json output error: %v", err)
	}
	var jsonResult DiagnoseReport
	if err = json.Unmarshal([]byte(jsonOutput), &jsonResult); err != nil {
		t.Fatalf("unmarshal json output: %v", err)
	}
	if !reflect.DeepEqual(jsonResult.Data, report.Data) {
		t.Fatalf("unexpected json result: %+v", jsonResult)
	}

	yamlOutput, err := FormatDiagnoseReport(report, OutputYAML)
	if err != nil {
		t.Fatalf("yaml output error: %v", err)
	}
	var yamlResult DiagnoseReport
	if err = yaml.Unmarshal([]byte(yamlOutput), &yamlResult); err != nil {
		t.Fatalf("unmarshal yaml output: %v", err)
	}
	if !reflect.DeepEqual(yamlResult.Data, report.Data) {
		t.Fatalf("unexpected yaml result: %+v", yamlResult)
	}
}

func TestCollectDiagnoseSnapshotStopsAfterStatusFailure(t *testing.T) {
	var transactionQueries, lockQueries int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case HealthCheckURL:
			fmt.Fprint(w, `{"code":"500","message":"status unavailable","success":false}`)
		case GlobalSessionQueryURL:
			transactionQueries++
		case GlobalLockQueryURL:
			lockQueries++
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	defer setTestAuth(t, server.URL, "test-token")()

	snapshot := CollectDiagnoseSnapshot(DiagnoseOptions{TCPTimeout: time.Second})
	if snapshot.Report == nil || snapshot.Report.Success == nil || *snapshot.Report.Success {
		t.Fatalf("unexpected report: %+v", snapshot.Report)
	}
	if transactionQueries != 0 || lockQueries != 0 {
		t.Fatalf("unexpected downstream queries: transactions=%d locks=%d", transactionQueries, lockQueries)
	}
	wantStages := []string{"config", "tcp connectivity", "login token", "status", "transaction query", "global lock query", "db/schema"}
	if len(snapshot.Report.Data) != len(wantStages) {
		t.Fatalf("len(report.Data) = %d, want %d", len(snapshot.Report.Data), len(wantStages))
	}
	for i, want := range wantStages {
		if snapshot.Report.Data[i].Stage != want {
			t.Fatalf("stage[%d] = %q, want %q", i, snapshot.Report.Data[i].Stage, want)
		}
	}
	if snapshot.Report.Data[4].Status != DiagnoseStatusSkip || snapshot.Report.Data[5].Status != DiagnoseStatusSkip {
		t.Fatalf("downstream stages = %+v, want skipped", snapshot.Report.Data[4:6])
	}
}

func TestCollectDiagnoseSnapshotKeepsIndependentQueriesAvailable(t *testing.T) {
	var transactionQueries, lockQueries int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case HealthCheckURL:
			fmt.Fprint(w, `{"code":"200","message":"success","success":true,"data":[]}`)
		case GlobalSessionQueryURL:
			transactionQueries++
			http.Error(w, "transaction query unavailable", http.StatusBadGateway)
		case GlobalLockQueryURL:
			lockQueries++
			fmt.Fprint(w, `{"code":"200","message":"success","success":true,"data":[]}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	defer setTestAuth(t, server.URL, "test-token")()

	snapshot := CollectDiagnoseSnapshot(DiagnoseOptions{TCPTimeout: time.Second})
	if transactionQueries != 1 || lockQueries != 1 {
		t.Fatalf("queries = transactions:%d locks:%d, want one of each", transactionQueries, lockQueries)
	}
	if snapshot.Transactions != nil || snapshot.TransactionErr == nil {
		t.Fatalf("transaction snapshot = %+v, error = %v, want query error", snapshot.Transactions, snapshot.TransactionErr)
	}
	if snapshot.Locks == nil || snapshot.LockErr != nil {
		t.Fatalf("lock snapshot = %+v, error = %v, want successful query", snapshot.Locks, snapshot.LockErr)
	}
	if snapshot.Report == nil || len(snapshot.Report.Data) != 7 {
		t.Fatalf("report = %+v, want all stages", snapshot.Report)
	}
	if snapshot.Report.Data[4].Status != DiagnoseStatusFail || snapshot.Report.Data[5].Status != DiagnoseStatusPass {
		t.Fatalf("query stages = %+v, want fail/pass", snapshot.Report.Data[4:6])
	}
}
