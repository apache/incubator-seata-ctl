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
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestLoginUsesJSONAndStoresToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var credentials map[string]string
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			t.Fatalf("decode login request: %v", err)
		}
		if credentials["username"] != `user"name` || credentials["password"] != `p\ss` {
			t.Fatalf("credentials = %#v", credentials)
		}
		fmt.Fprint(w, `{"code":"200","message":"success","data":"new-token","success":true}`)
	}))
	defer server.Close()

	auth := GetAuth()
	old := *auth
	defer func() { *auth = old }()
	setAuthAddress(t, auth, server.URL)
	auth.Username = `user"name`
	auth.Password = `p\ss`
	if err := auth.Login(); err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if token, err := auth.GetToken(); err != nil || token != "new-token" {
		t.Fatalf("token = %q, err = %v", token, err)
	}
}

func TestLoginRejectsInvalidResponsesAndClearsToken(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		want       string
	}{
		{name: "http error", statusCode: http.StatusBadGateway, body: "upstream failed", want: "http status 502"},
		{name: "business error", statusCode: http.StatusOK, body: `{"code":"401","message":"unauthorized"}`, want: "unauthorized"},
		{name: "empty token", statusCode: http.StatusOK, body: `{"code":"200","message":"success","data":""}`, want: "empty token"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(test.statusCode)
				fmt.Fprint(w, test.body)
			}))
			defer server.Close()

			auth := GetAuth()
			old := *auth
			defer func() { *auth = old }()
			setAuthAddress(t, auth, server.URL)
			auth.token = "old-token"

			err := auth.Login()
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Login error = %v, want %q", err, test.want)
			}
			if _, err = auth.GetToken(); err == nil {
				t.Fatal("GetToken succeeded after failed login")
			}
		})
	}
}

func setAuthAddress(t *testing.T, auth *Auth, serverURL string) {
	t.Helper()
	parsed, err := url.Parse(serverURL)
	if err != nil {
		t.Fatalf("parse server URL: %v", err)
	}
	host, portString, err := net.SplitHostPort(parsed.Host)
	if err != nil {
		t.Fatalf("split server address: %v", err)
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		t.Fatalf("parse server port: %v", err)
	}
	auth.ServerIP = host
	auth.ServerPort = port
}
