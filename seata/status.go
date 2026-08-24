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
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type NodeStatusResponse struct {
	BaseResponse
	Data []NodeStatus `json:"data" yaml:"data"`
}

type NodeStatus struct {
	Address string `json:"address" yaml:"address"`
	Status  string `json:"status" yaml:"status"`
	Type    string `json:"type" yaml:"type"`
}

func QueryStatus() (*NodeStatusResponse, error) {
	return queryStatus(NewConsoleClient())
}

func queryStatus(client *ConsoleClient) (*NodeStatusResponse, error) {
	var response NodeStatusResponse
	if err := client.Get(HealthCheckURL, nil, &response); err != nil {
		return nil, err
	}
	if err := checkConsoleCode(response.Code, response.Message, "query status failed"); err != nil {
		return nil, err
	}
	return &response, nil
}

func GetStatus() {
	response, err := QueryStatus()
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := FormatNodeStatusResponse(response, OutputTable)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintln(os.Stdout, result)
}

func FormatNodeStatusResponse(response *NodeStatusResponse, output string) (string, error) {
	output, err := NormalizeOutput(output)
	if err != nil {
		return "", err
	}
	switch output {
	case OutputTable:
		return FormatNodeStatusTable(response.Data), nil
	default:
		return FormatStructuredOutput(response, output)
	}
}

func FormatNodeStatusTable(statuses []NodeStatus) string {
	t := table.NewWriter()
	t.Style().Format.Header = text.FormatDefault
	t.AppendHeader(table.Row{"type", "address", "status"})
	for _, status := range statuses {
		t.AppendRow(table.Row{status.Type, status.Address, status.Status})
	}
	return t.Render()
}
