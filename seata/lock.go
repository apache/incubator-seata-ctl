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
	"errors"
	"net/url"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type GlobalLockQuery struct {
	XID           string
	TableName     string
	TransactionID string
	BranchID      string
	PK            string
	ResourceID    string
	PageNum       int
	PageSize      int
	TimeStart     *int64
	TimeEnd       *int64
}

type GlobalLockPageResult struct {
	Code     string       `json:"code" yaml:"code"`
	Message  string       `json:"message" yaml:"message"`
	Success  *bool        `json:"success" yaml:"success"`
	PageSize *int         `json:"pageSize" yaml:"pageSize"`
	PageNum  *int         `json:"pageNum" yaml:"pageNum"`
	CurrPage *int         `json:"currPage" yaml:"currPage"`
	Total    *int         `json:"total" yaml:"total"`
	Pages    *int         `json:"pages" yaml:"pages"`
	Data     []GlobalLock `json:"data" yaml:"data"`
}

type GlobalLock struct {
	XID           *string `json:"xid" yaml:"xid"`
	TransactionID *string `json:"transactionId" yaml:"transactionId"`
	BranchID      *string `json:"branchId" yaml:"branchId"`
	ResourceID    *string `json:"resourceId" yaml:"resourceId"`
	TableName     *string `json:"tableName" yaml:"tableName"`
	PK            *string `json:"pk" yaml:"pk"`
	RowKey        *string `json:"rowKey" yaml:"rowKey"`
	Vgroup        *string `json:"vgroup" yaml:"vgroup"`
	GmtCreate     *int64  `json:"gmtCreate" yaml:"gmtCreate"`
	GmtModified   *int64  `json:"gmtModified" yaml:"gmtModified"`
}

type GlobalLockCheckResult struct {
	Code    string `json:"code" yaml:"code"`
	Message string `json:"message" yaml:"message"`
	Success *bool  `json:"success" yaml:"success"`
	Data    *bool  `json:"data" yaml:"data"`
}

func QueryGlobalLocks(query GlobalLockQuery) (*GlobalLockPageResult, error) {
	return queryGlobalLocks(NewConsoleClient(), query)
}

func queryGlobalLocks(client *ConsoleClient, query GlobalLockQuery) (*GlobalLockPageResult, error) {
	if query.PageNum <= 0 {
		return nil, errors.New("page-num must be greater than 0")
	}
	if query.PageSize <= 0 {
		return nil, errors.New("page-size must be greater than 0")
	}

	var response GlobalLockPageResult
	if err := client.Get(GlobalLockQueryURL, query.values(), &response); err != nil {
		return nil, err
	}
	if err := checkConsoleCode(response.Code, response.Message, "query global locks failed"); err != nil {
		return nil, err
	}
	return &response, nil
}

func CheckGlobalLock(xid string, branchID string) (*GlobalLockCheckResult, error) {
	return checkGlobalLock(NewConsoleClient(), xid, branchID)
}

func checkGlobalLock(client *ConsoleClient, xid string, branchID string) (*GlobalLockCheckResult, error) {
	if xid == "" {
		return nil, errors.New("xid is required")
	}
	if branchID == "" {
		return nil, errors.New("branch-id is required")
	}

	values := url.Values{}
	values.Set("xid", xid)
	values.Set("branchId", branchID)
	var response GlobalLockCheckResult
	if err := client.Get(GlobalLockCheckURL, values, &response); err != nil {
		return nil, err
	}
	if err := checkConsoleCode(response.Code, response.Message, "check global lock failed"); err != nil {
		return nil, err
	}
	return &response, nil
}

func (query GlobalLockQuery) values() url.Values {
	values := url.Values{}
	values.Set("pageNum", strconv.Itoa(query.PageNum))
	values.Set("pageSize", strconv.Itoa(query.PageSize))
	if query.XID != "" {
		values.Set("xid", query.XID)
	}
	if query.TableName != "" {
		values.Set("tableName", query.TableName)
	}
	if query.TransactionID != "" {
		values.Set("transactionId", query.TransactionID)
	}
	if query.BranchID != "" {
		values.Set("branchId", query.BranchID)
	}
	if query.PK != "" {
		values.Set("pk", query.PK)
	}
	if query.ResourceID != "" {
		values.Set("resourceId", query.ResourceID)
	}
	if query.TimeStart != nil {
		values.Set("timeStart", strconv.FormatInt(*query.TimeStart, 10))
	}
	if query.TimeEnd != nil {
		values.Set("timeEnd", strconv.FormatInt(*query.TimeEnd, 10))
	}
	return values
}

func FormatGlobalLockPage(response *GlobalLockPageResult, output string) (string, error) {
	output, err := NormalizeOutput(output)
	if err != nil {
		return "", err
	}
	switch output {
	case OutputTable:
		return FormatGlobalLockTable(response.Data), nil
	default:
		return FormatStructuredOutput(response, output)
	}
}

func FormatGlobalLockCheck(response *GlobalLockCheckResult, output string) (string, error) {
	output, err := NormalizeOutput(output)
	if err != nil {
		return "", err
	}
	switch output {
	case OutputTable:
		locked := ""
		if response.Data != nil {
			locked = strconv.FormatBool(*response.Data)
		}
		t := table.NewWriter()
		t.Style().Format.Header = text.FormatDefault
		t.AppendHeader(table.Row{"locked", "message"})
		t.AppendRow(table.Row{locked, response.Message})
		return t.Render(), nil
	default:
		return FormatStructuredOutput(response, output)
	}
}

func FormatGlobalLockTable(locks []GlobalLock) string {
	t := table.NewWriter()
	t.Style().Format.Header = text.FormatDefault
	t.AppendHeader(table.Row{"xid", "transaction_id", "branch_id", "resource_id", "table_name", "pk", "row_key", "vgroup"})
	for _, lock := range locks {
		t.AppendRow(table.Row{
			stringValue(lock.XID),
			stringValue(lock.TransactionID),
			stringValue(lock.BranchID),
			stringValue(lock.ResourceID),
			stringValue(lock.TableName),
			stringValue(lock.PK),
			stringValue(lock.RowKey),
			stringValue(lock.Vgroup),
		})
	}
	return t.Render()
}
