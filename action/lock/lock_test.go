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

package lock

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
	"testing"

	"github.com/seata/seata-ctl/seata"
	"github.com/spf13/cobra"
)

func TestListCommandBuildsQuery(t *testing.T) {
	server := newConsoleServer(t, func(r *http.Request) {
		wantQuery := map[string]string{
			"xid":           "xid-1",
			"tableName":     "account",
			"transactionId": "1001",
			"branchId":      "2001",
			"pk":            "1",
			"resourceId":    "jdbc:mysql://127.0.0.1:3306/seata",
			"pageNum":       "2",
			"pageSize":      "5",
			"timeStart":     "1710000000000",
			"timeEnd":       "1710003600000",
		}
		for key, want := range wantQuery {
			if got := r.URL.Query().Get(key); got != want {
				t.Fatalf("query[%s] = %q, want %q", key, got, want)
			}
		}
	})
	defer server.Close()
	loginToServer(t, server)

	cmd := newListCommand()
	output, err := executeCommand(t, cmd,
		"--xid", "xid-1",
		"--table-name", "account",
		"--transaction-id", "1001",
		"--branch-id", "2001",
		"--pk", "1",
		"--resource-id", "jdbc:mysql://127.0.0.1:3306/seata",
		"--page-num", "2",
		"--page-size", "5",
		"--time-start", "1710000000000",
		"--time-end", "1710003600000",
		"--output", "json",
	)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}
	if !strings.Contains(output, `"pageNum": 2`) || !strings.Contains(output, `"data"`) {
		t.Fatalf("output = %s, want JSON page output", output)
	}
}

func TestCheckCommandBuildsQuery(t *testing.T) {
	server := newConsoleServer(t, nil)
	defer server.Close()
	loginToServer(t, server)

	cmd := newCheckCommand()
	output, err := executeCommand(t, cmd, "--xid", "xid-1", "--branch-id", "2001")
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}
	if !strings.Contains(output, "locked") && !strings.Contains(output, "true") {
		t.Fatalf("output = %s, want lock status", output)
	}
}

func newConsoleServer(t *testing.T, onQuery func(*http.Request)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case seata.LoginURL:
			fmt.Fprint(w, `{"code":"200","message":"success","data":"test-token","success":true}`)
		case seata.GlobalLockQueryURL:
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
			fmt.Fprintf(w, `{"code":"200","message":"success","success":true,"pageSize":%s,"pageNum":%s,"total":1,"pages":1,"data":[{"xid":"xid-1","transactionId":"1001","branchId":"2001","resourceId":"jdbc:mysql://127.0.0.1:3306/seata","tableName":"account","pk":"1","rowKey":"1","vgroup":"default_tx_group"}]}`, pageSize, pageNum)
		case seata.GlobalLockCheckURL:
			fmt.Fprint(w, `{"code":"200","message":"success","success":true,"data":true}`)
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

