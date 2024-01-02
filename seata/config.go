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
	"os"
)

type KVConfigResponse struct {
	BaseResponse
	Data map[string]string `json:"data"`
}

type DiffConfigResponse struct {
	BaseResponse
	Data map[string][]string `json:"data"`
}

func GetConfigurations(params []string) (string, error) {
	urlStr := HTTPProtocol + GetAuth().GetAddress() + GetConfigurationURL
	request, err := BuildPostRequestWithArrayData(urlStr, params)
	if err != nil {
		return "", err
	}

	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response KVConfigResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Code != CodeOK {
		return "", errors.New(response.Message)
	}

	return FormatKVResponse(response.Data), nil
}

func SetConfiguration(data map[string]string, configType ConfigType) (string, error) {
	urlStr := HTTPProtocol + GetAuth().GetAddress()
	switch configType {
	case RegistryConf:
		urlStr = urlStr + RegistryConfigurationURL
	case ConfigCenterConf:
		urlStr = urlStr + ConfigCenterConfigurationURL
	default:
		urlStr = urlStr + ConfigurationURL
	}
	request, err := BuildPostRequestWithMapData(urlStr, data)
	if err != nil {
		return "", err
	}

	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response DiffConfigResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Code != CodeOK {
		return "", errors.New(response.Message)
	}

	return FormatDiffResponse(response.Data), nil
}

func ReloadConfiguration() {
	url := HTTPProtocol + GetAuth().GetAddress() + ReloadConfigurationURL
	token, err := GetAuth().GetToken()
	if err != nil {
		fmt.Println("Please login again!")
		os.Exit(0)
	}

	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Set("authorization", token)
	request.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var response BaseResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
	}

	if response.Code != CodeOK {
		fmt.Println(response.Message)
	} else {
		fmt.Println("Reload Successful!")
	}
}
