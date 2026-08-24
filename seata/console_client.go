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
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ConsoleClient struct {
	HTTPClient *http.Client
}

func NewConsoleClient() *ConsoleClient {
	return &ConsoleClient{HTTPClient: defaultHTTPClient}
}

func (client *ConsoleClient) Get(path string, values url.Values, target interface{}) error {
	token, err := GetAuth().GetToken()
	if err != nil {
		return errors.New("please login")
	}

	request, err := http.NewRequest(http.MethodGet, consoleURL(path, values), nil)
	if err != nil {
		return err
	}
	request.Header.Set("authorization", token)

	httpClient := client.HTTPClient
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("get %s: %w", path, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("get %s: http status %d: %s", path, response.StatusCode, strings.TrimSpace(string(body)))
	}
	if err = json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode %s response: %w", path, err)
	}
	return nil
}

func consoleURL(path string, values url.Values) string {
	rawURL := HTTPProtocol + GetAuth().GetAddress() + path
	if len(values) == 0 {
		return rawURL
	}
	return rawURL + "?" + values.Encode()
}

func checkConsoleCode(code string, message string, fallback string) error {
	if code == CodeOK {
		return nil
	}
	if message == "" {
		message = fallback
	}
	return errors.New(message)
}
