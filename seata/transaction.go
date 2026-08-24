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
	"fmt"
	"net/url"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type GlobalSessionQuery struct {
	XID             string
	ApplicationID   string
	Status          *int
	TransactionName string
	Vgroup          string
	WithBranch      *bool
	PageNum         int
	PageSize        int
	TimeStart       *int64
	TimeEnd         *int64
}

type GlobalSessionPageResult struct {
	Code     string          `json:"code" yaml:"code"`
	Message  string          `json:"message" yaml:"message"`
	Success  *bool           `json:"success" yaml:"success"`
	PageSize *int            `json:"pageSize" yaml:"pageSize"`
	PageNum  *int            `json:"pageNum" yaml:"pageNum"`
	CurrPage *int            `json:"currPage" yaml:"currPage"`
	Total    *int            `json:"total" yaml:"total"`
	Pages    *int            `json:"pages" yaml:"pages"`
	Data     []GlobalSession `json:"data" yaml:"data"`
}

type GlobalSession struct {
	XID                     *string          `json:"xid" yaml:"xid"`
	TransactionID           *string          `json:"transactionId" yaml:"transactionId"`
	Status                  *int             `json:"status" yaml:"status"`
	ApplicationID           *string          `json:"applicationId" yaml:"applicationId"`
	TransactionServiceGroup *string          `json:"transactionServiceGroup" yaml:"transactionServiceGroup"`
	TransactionName         *string          `json:"transactionName" yaml:"transactionName"`
	Timeout                 *int64           `json:"timeout" yaml:"timeout"`
	BeginTime               *int64           `json:"beginTime" yaml:"beginTime"`
	ApplicationData         *string          `json:"applicationData" yaml:"applicationData"`
	GmtCreate               *int64           `json:"gmtCreate" yaml:"gmtCreate"`
	GmtModified             *int64           `json:"gmtModified" yaml:"gmtModified"`
	BranchSessionVOs        *[]BranchSession `json:"branchSessionVOs" yaml:"branchSessionVOs"`
}

type BranchSession struct {
	XID             *string `json:"xid" yaml:"xid"`
	TransactionID   *string `json:"transactionId" yaml:"transactionId"`
	BranchID        *string `json:"branchId" yaml:"branchId"`
	ResourceGroupID *string `json:"resourceGroupId" yaml:"resourceGroupId"`
	ResourceID      *string `json:"resourceId" yaml:"resourceId"`
	BranchType      *string `json:"branchType" yaml:"branchType"`
	Status          *int    `json:"status" yaml:"status"`
	ClientID        *string `json:"clientId" yaml:"clientId"`
	ApplicationData *string `json:"applicationData" yaml:"applicationData"`
	GmtCreate       *int64  `json:"gmtCreate" yaml:"gmtCreate"`
	GmtModified     *int64  `json:"gmtModified" yaml:"gmtModified"`
}

func QueryGlobalSessions(query GlobalSessionQuery) (*GlobalSessionPageResult, error) {
	return queryGlobalSessions(NewConsoleClient(), query)
}

func queryGlobalSessions(client *ConsoleClient, query GlobalSessionQuery) (*GlobalSessionPageResult, error) {
	if query.PageNum <= 0 {
		return nil, errors.New("page-num must be greater than 0")
	}
	if query.PageSize <= 0 {
		return nil, errors.New("page-size must be greater than 0")
	}

	var response GlobalSessionPageResult
	if err := client.Get(GlobalSessionQueryURL, query.values(), &response); err != nil {
		return nil, err
	}
	if err := checkConsoleCode(response.Code, response.Message, "query global sessions failed"); err != nil {
		return nil, err
	}

	return &response, nil
}

func QueryGlobalSessionByXID(xid string) (*GlobalSession, error) {
	withBranch := true
	response, err := QueryGlobalSessions(GlobalSessionQuery{
		XID:        xid,
		WithBranch: &withBranch,
		PageNum:    1,
		PageSize:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("global session %q not found", xid)
	}
	return &response.Data[0], nil
}

func (query GlobalSessionQuery) values() url.Values {
	values := url.Values{}
	values.Set("pageNum", strconv.Itoa(query.PageNum))
	values.Set("pageSize", strconv.Itoa(query.PageSize))
	if query.XID != "" {
		values.Set("xid", query.XID)
	}
	if query.ApplicationID != "" {
		values.Set("applicationId", query.ApplicationID)
	}
	if query.Status != nil {
		values.Set("status", strconv.Itoa(*query.Status))
	}
	if query.TransactionName != "" {
		values.Set("transactionName", query.TransactionName)
	}
	if query.Vgroup != "" {
		values.Set("vgroup", query.Vgroup)
	}
	if query.WithBranch != nil {
		values.Set("withBranch", strconv.FormatBool(*query.WithBranch))
	}
	if query.TimeStart != nil {
		values.Set("timeStart", strconv.FormatInt(*query.TimeStart, 10))
	}
	if query.TimeEnd != nil {
		values.Set("timeEnd", strconv.FormatInt(*query.TimeEnd, 10))
	}
	return values
}

func FormatGlobalSessionPage(response *GlobalSessionPageResult, output string) (string, error) {
	output, err := NormalizeOutput(output)
	if err != nil {
		return "", err
	}
	switch output {
	case OutputTable:
		return FormatGlobalSessionTable(response.Data), nil
	default:
		return FormatStructuredOutput(response, output)
	}
}

func FormatGlobalSessionDetail(session *GlobalSession, output string) (string, error) {
	output, err := NormalizeOutput(output)
	if err != nil {
		return "", err
	}
	switch output {
	case OutputTable:
		return FormatGlobalSessionDetailTable(session), nil
	default:
		return FormatStructuredOutput(session, output)
	}
}

func FormatGlobalSessionTable(sessions []GlobalSession) string {
	t := table.NewWriter()
	t.Style().Format.Header = text.FormatDefault
	t.AppendHeader(table.Row{"xid", "transaction_id", "status", "application_id", "vgroup", "transaction_name", "begin_time", "timeout"})
	for _, session := range sessions {
		t.AppendRow(table.Row{
			stringValue(session.XID),
			stringValue(session.TransactionID),
			intValue(session.Status),
			stringValue(session.ApplicationID),
			stringValue(session.TransactionServiceGroup),
			stringValue(session.TransactionName),
			int64Value(session.BeginTime),
			int64Value(session.Timeout),
		})
	}
	return t.Render()
}

func FormatGlobalSessionDetailTable(session *GlobalSession) string {
	t := table.NewWriter()
	t.Style().Format.Header = text.FormatDefault
	t.AppendHeader(table.Row{"field", "value"})
	t.AppendRows([]table.Row{
		{"xid", stringValue(session.XID)},
		{"transaction_id", stringValue(session.TransactionID)},
		{"status", intValue(session.Status)},
		{"application_id", stringValue(session.ApplicationID)},
		{"vgroup", stringValue(session.TransactionServiceGroup)},
		{"transaction_name", stringValue(session.TransactionName)},
		{"begin_time", int64Value(session.BeginTime)},
		{"timeout", int64Value(session.Timeout)},
		{"application_data", stringValue(session.ApplicationData)},
		{"gmt_create", int64Value(session.GmtCreate)},
		{"gmt_modified", int64Value(session.GmtModified)},
	})

	if session.BranchSessionVOs == nil || len(*session.BranchSessionVOs) == 0 {
		return t.Render()
	}

	branches := table.NewWriter()
	branches.Style().Format.Header = text.FormatDefault
	branches.AppendHeader(table.Row{"branch_id", "status", "resource_id", "branch_type", "client_id", "gmt_create", "gmt_modified"})
	for _, branch := range *session.BranchSessionVOs {
		branches.AppendRow(table.Row{
			stringValue(branch.BranchID),
			intValue(branch.Status),
			stringValue(branch.ResourceID),
			stringValue(branch.BranchType),
			stringValue(branch.ClientID),
			int64Value(branch.GmtCreate),
			int64Value(branch.GmtModified),
		})
	}
	return t.Render() + "\n" + branches.Render()
}
