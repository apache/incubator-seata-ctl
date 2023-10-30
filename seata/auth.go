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
)

var auth Auth

type Auth struct {
	ServerIp   string
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
	return auth.ServerIp + ":" + strconv.Itoa(auth.ServerPort)
}

func GetAuth() *Auth {
	return &auth
}

func (auth *Auth) Login() error {
	url := HTTPProtocol + auth.GetAddress() + LoginURL
	jsonStr := []byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, auth.Username, auth.Password))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var jsonResp Response
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return err
	}
	auth.token = jsonResp.Data
	return nil
}
