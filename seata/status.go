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
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"net/http"
	"os"
)

type NodeStatusResponse struct {
	BaseResponse
	Data []NodeStatus `json:"data"`
}

type NodeStatus struct {
	Address string `json:"address"`
	Status  string `json:"status"`
	Type    string `json:"type"`
}

func GetStatus() {
	url := HTTPProtocol + GetAuth().GetAddress() + HealthCheckURL
	token, err := GetAuth().GetToken()
	if err != nil {
		fmt.Println("Please login again!")
		os.Exit(0)
	}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", token)
	resp, err := (&http.Client{}).Do(request)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var response NodeStatusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
	}

	if response.Code != "200" {
		fmt.Println(response.Message)
	}

	t := table.NewWriter()
	header := table.Row{"type", "address", "status"}
	t.AppendHeader(header)
	for _, data := range response.Data {
		row := table.Row{data.Type, data.Address, data.Status}
		t.AppendRow(row)
	}
	fmt.Println(t.Render())
	t.Style()
}
