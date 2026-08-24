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

package diagnose

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

func TestRunCommandBuildsReport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case seata.LoginURL:
			fmt.Fprint(w, `{"code":"200","message":"success","data":"test-token","success":true}`)
		case seata.HealthCheckURL:
			if got := r.Header.Get("authorization"); got != "test-token" {
				t.Fatalf("authorization = %q, want test-token", got)
			}
			fmt.Fprint(w, `{"code":"200","message":"success","success":true,"data":[{"type":"nacos","address":"127.0.0.1:7091","status":"ok"}]}`)
		case seata.GlobalSessionQueryURL:
			if got := r.Header.Get("authorization"); got != "test-token" {
				t.Fatalf("authorization = %q, want test-token", got)
			}
			fmt.Fprint(w, `{"code":"200","message":"success","success":true,"pageSize":1,"pageNum":1,"total":1,"pages":1,"data":[{"xid":"xid-1","transactionId":"1001","status":1,"applicationId":"order-service","transactionServiceGroup":"default_tx_group","transactionName":"createOrder","timeout":60000,"beginTime":1710000000000}]}`)
		case seata.GlobalLockQueryURL:
			if got := r.Header.Get("authorization"); got != "test-token" {
				t.Fatalf("authorization = %q, want test-token", got)
			}
			fmt.Fprint(w, `{"code":"200","message":"success","success":true,"pageSize":1,"pageNum":1,"total":1,"pages":1,"data":[{"xid":"xid-1","transactionId":"1001","branchId":"2001","resourceId":"jdbc:mysql://127.0.0.1:3306/seata","tableName":"account","pk":"1","rowKey":"1","vgroup":"default_tx_group"}]}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	loginToServer(t, server)

	cmd := newRunCommand()
	output, err := executeCommand(t, cmd, "--output", "json")
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}
	if !strings.Contains(output, `"stage": "status"`) || !strings.Contains(output, `"stage": "transaction query"`) || !strings.Contains(output, `"stage": "global lock query"`) {
		t.Fatalf("output = %s, want JSON diagnose report", output)
	}
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
