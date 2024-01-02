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
	"net/http"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
)

func BuildPostRequestWithArrayData(urlStr string, data []string) (*http.Request, error) {
	token, err := GetAuth().GetToken()
	if err != nil {
		return nil, errors.New("please login")
	}

	body, _ := json.Marshal(data)

	request, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(body))
	request.Header.Set("authorization", token)
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

func BuildPostRequestWithMapData(urlStr string, data map[string]string) (*http.Request, error) {
	token, err := GetAuth().GetToken()
	if err != nil {
		return nil, errors.New("please login")
	}

	body, _ := json.Marshal(data)

	request, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(body))
	request.Header.Set("authorization", token)
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

func FormatKVResponse(kv map[string]string) string {
	t := table.NewWriter()
	header := table.Row{"key", "value"}
	t.AppendHeader(header)

	// Make output in order
	var keys []string
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		row := table.Row{key, kv[key]}
		t.AppendRow(row)
	}
	return t.Render()
}

func FormatDiffResponse(kv map[string][]string) string {
	t := table.NewWriter()
	header := table.Row{"key", "from", "to"}
	t.AppendHeader(header)

	// Make output in order
	var keys []string
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		row := table.Row{key, kv[key][0], kv[key][1]}
		t.AppendRow(row)
	}
	return t.Render()
}
