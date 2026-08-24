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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var auth Auth

type Auth struct {
	ServerIP   string
	ServerPort int
	Username   string
	Password   string
	token      string
}

type Response struct {
	Code    string
	Message string
	Data    string
	Success bool
}

func (auth *Auth) GetToken() (string, error) {
	if auth.token == "" {
		return auth.token, errors.New("login failed")
	}
	return auth.token, nil
}

func (auth *Auth) GetAddress() string {
	return auth.ServerIP + ":" + strconv.Itoa(auth.ServerPort)
}

func GetAuth() *Auth {
	return &auth
}

func (auth *Auth) Login() error {
	auth.token = ""
	url := HTTPProtocol + auth.GetAddress() + LoginURL
	jsonStr, err := json.Marshal(map[string]string{
		"username": auth.Username,
		"password": auth.Password,
	})
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := defaultHTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("login failed: http status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var jsonResp Response
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return err
	}
	if jsonResp.Code != CodeOK {
		if jsonResp.Message == "" {
			return errors.New("login failed")
		}
		return errors.New(jsonResp.Message)
	}
	if jsonResp.Data == "" {
		return errors.New("login failed: empty token")
	}
	auth.token = jsonResp.Data
	return nil
}
