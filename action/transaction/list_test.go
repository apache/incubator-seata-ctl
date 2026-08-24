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

package transaction

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/seata/seata-ctl/seata"
	"github.com/spf13/cobra"
)

func TestListCommandBuildsQuery(t *testing.T) {
	var (
		mu        sync.Mutex
		gotHeader string
		gotQuery  url.Values
	)
	server := newConsoleServer(t, func(r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		gotHeader = r.Header.Get("authorization")
		gotQuery = r.URL.Query()
	})
	defer server.Close()
	loginToServer(t, server)

	cmd := newListCommand()
	output, err := executeCommand(t, cmd,
		"--xid", "xid-1",
		"--application-id", "order-service",
		"--status", "1",
		"--transaction-name", "createOrder",
		"--vgroup", "default_tx_group",
		"--with-branch",
		"--page-num", "2",
		"--page-size", "5",
		"--time-start", "1710000000000",
		"--time-end", "1710003600000",
		"--output", "json",
	)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}
	if !strings.Contains(output, `"pageNum": 2`) {
		t.Fatalf("output = %s, want JSON pageNum", output)
	}

	mu.Lock()
	defer mu.Unlock()
	if gotHeader != "test-token" {
		t.Fatalf("authorization = %q, want test-token", gotHeader)
	}
	wantQuery := map[string]string{
		"xid":             "xid-1",
		"applicationId":   "order-service",
		"status":          "1",
		"transactionName": "createOrder",
		"vgroup":          "default_tx_group",
		"withBranch":      "true",
		"pageNum":         "2",
		"pageSize":        "5",
		"timeStart":       "1710000000000",
		"timeEnd":         "1710003600000",
	}
	for key, want := range wantQuery {
		if got := gotQuery.Get(key); got != want {
			t.Fatalf("query[%s] = %q, want %q", key, got, want)
		}
	}
}

func TestListCommandResetsFlagsAfterParseError(t *testing.T) {
	server := newConsoleServer(t, nil)
	defer server.Close()
	loginToServer(t, server)

	cmd := newListCommand()
	_, err := executeCommand(t, cmd, "--output", "json", "--unknown")
	if err == nil || !strings.Contains(err.Error(), "unknown flag") {
		t.Fatalf("error = %v, want unknown flag", err)
	}

	output, err := executeCommand(t, cmd)
	if err != nil {
		t.Fatalf("execute command after parse error: %v", err)
	}
	if strings.Contains(output, `"code":`) {
		t.Fatalf("output = %s, want default table output after parse error", output)
	}
	if !strings.Contains(output, "transaction_id") {
		t.Fatalf("output = %s, want table header", output)
	}
}

func TestListCommandResetsFlagsAfterHelp(t *testing.T) {
	server := newConsoleServer(t, nil)
	defer server.Close()
	loginToServer(t, server)

	cmd := newListCommand()
	helpOutput, err := executeCommand(t, cmd, "--output", "json", "--help")
	if err != nil {
		t.Fatalf("execute help command: %v", err)
	}
	if !strings.Contains(helpOutput, "List global transactions") {
		t.Fatalf("help output = %s, want transaction list help", helpOutput)
	}

	output, err := executeCommand(t, cmd)
	if err != nil {
		t.Fatalf("execute command after help: %v", err)
	}
	if strings.Contains(output, `"code":`) {
		t.Fatalf("output = %s, want default table output after help", output)
	}
	if !strings.Contains(output, "transaction_id") {
		t.Fatalf("output = %s, want table header", output)
	}
}

func newConsoleServer(t *testing.T, onQuery func(*http.Request)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case seata.LoginURL:
			fmt.Fprint(w, `{"code":"200","message":"success","data":"test-token","success":true}`)
		case seata.GlobalSessionQueryURL:
			if onQuery != nil {
				onQuery(r)
			}
			pageNum := r.URL.Query().Get("pageNum")
			if pageNum == "" {
				pageNum = "1"
			}
			pageSize := r.URL.Query().Get("pageSize")
			if pageSize == "" {
				pageSize = "20"
			}
			fmt.Fprintf(w, `{"code":"200","message":"success","success":true,"pageSize":%s,"pageNum":%s,"total":1,"pages":1,"data":[{"xid":"xid-1","transactionId":"1001","status":1,"applicationId":"order-service","transactionServiceGroup":"default_tx_group","transactionName":"createOrder","timeout":60000,"beginTime":1710000000000}]}`, pageSize, pageNum)
		default:
			http.NotFound(w, r)
		}
	}))
}

func loginToServer(t *testing.T, server *httptest.Server) {
	t.Helper()
	parsedURL, err := url.Parse(server.URL)
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
	auth := seata.GetAuth()
	auth.ServerIP = host
	auth.ServerPort = port
	auth.Username = "seata"
	auth.Password = "seata"
	if err = auth.Login(); err != nil {
		t.Fatalf("login to test server: %v", err)
	}
}

func executeCommand(t *testing.T, cmd *cobra.Command, args ...string) (string, error) {
	t.Helper()
	cmd.SetArgs(args)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	var output bytes.Buffer
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = writer
	execErr := cmd.Execute()
	_ = writer.Close()
	os.Stdout = oldStdout
	if _, err = io.Copy(&output, reader); err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	_ = reader.Close()
	return output.String(), execErr
}
