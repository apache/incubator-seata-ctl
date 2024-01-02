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
	"io"
	"net/http"
	"os"
	"strconv"
)

type TxnResponse struct {
	BaseResponse
	Data string
}

func BeginTxn(timeout int) {
	url := HTTPProtocol + GetAuth().GetAddress() + TryBeginTxnURL
	url = url + "?timeout=" + strconv.Itoa(timeout)
	token, err := GetAuth().GetToken()
	if err != nil {
		fmt.Println("Please login again!")
		os.Exit(0)
	}
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Set("authorization", token)
	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var response TxnResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
	}

	if response.Code != CodeOK {
		fmt.Println(response.Message)
	} else {
		fmt.Printf("Try an example txn successfully, xid=%s\n", response.Data)
	}
}

func CommitTxn(xid string) {
	url := HTTPProtocol + GetAuth().GetAddress() + TryCommitTxnURL
	url = url + "?xid=" + xid
	token, err := GetAuth().GetToken()
	if err != nil {
		fmt.Println("Please login again!")
		os.Exit(0)
	}
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Set("authorization", token)
	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var response TxnResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
	}

	if response.Code != CodeOK {
		fmt.Println(response.Message)
	} else {
		fmt.Printf("Commit txn successfully, xid=%s\n", response.Data)
	}
}

func RollbackTxn(xid string) {
	url := HTTPProtocol + GetAuth().GetAddress() + TryRollBackTxnURL
	url = url + "?xid=" + xid
	token, err := GetAuth().GetToken()
	if err != nil {
		fmt.Println("Please login again!")
		os.Exit(0)
	}
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Set("authorization", token)
	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var response TxnResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
	}

	if response.Code != CodeOK {
		fmt.Println(response.Message)
	} else {
		fmt.Printf("Rollback txn successfully, xid=%s\n", response.Data)
	}
}
