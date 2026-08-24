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
	"net"
	"strings"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

const (
	DiagnoseStatusPass = "PASS"
	DiagnoseStatusWarn = "WARN"
	DiagnoseStatusFail = "FAIL"
	DiagnoseStatusSkip = "SKIP"
)

type DiagnoseOptions struct {
	CheckDB      bool
	DBAddress    string
	TCPTimeout   time.Duration
	PageNum      int
	PageSize     int
	LockPageNum  int
	LockPageSize int
}

type DiagnoseReport struct {
	Code    string        `json:"code" yaml:"code"`
	Message string        `json:"message" yaml:"message"`
	Success *bool         `json:"success" yaml:"success"`
	Data    []CheckResult `json:"data" yaml:"data"`
}

type DiagnoseSnapshot struct {
	Report         *DiagnoseReport
	Transactions   *GlobalSessionPageResult
	TransactionErr error
	Locks          *GlobalLockPageResult
	LockErr        error
}

type CheckResult struct {
	Stage    string `json:"stage" yaml:"stage"`
	Status   string `json:"status" yaml:"status"`
	Message  string `json:"message" yaml:"message"`
	Detail   string `json:"detail,omitempty" yaml:"detail,omitempty"`
	Evidence string `json:"evidence,omitempty" yaml:"evidence,omitempty"`
}

func RunDiagnostics(opts DiagnoseOptions) (*DiagnoseReport, error) {
	snapshot := CollectDiagnoseSnapshot(opts)
	return snapshot.Report, nil
}

func CollectDiagnoseSnapshot(opts DiagnoseOptions) *DiagnoseSnapshot {
	opts = normalizeDiagnoseOptions(opts)
	snapshot := &DiagnoseSnapshot{}
	results := make([]CheckResult, 0, 7)
	passed := true
	add := func(result CheckResult) {
		if result.Status == DiagnoseStatusFail {
			passed = false
		}
		results = append(results, result)
	}
	finish := func() {
		snapshot.Report = &DiagnoseReport{
			Code:    CodeOK,
			Message: "diagnostics finished",
			Success: boolPtr(passed),
			Data:    results,
		}
	}

	address := strings.TrimSpace(GetAuth().GetAddress())
	if address == "" || address == ":0" || address == ":" {
		add(CheckResult{Stage: "config", Status: DiagnoseStatusFail, Message: "server address is not configured"})
		addSkipped(add, "config failed", "tcp connectivity", "login token", "status", "transaction query", "global lock query")
		add(diagnoseDBResult(opts, DiagnoseStatusSkip, "skipped because config failed"))
		finish()
		return snapshot
	}
	add(CheckResult{Stage: "config", Status: DiagnoseStatusPass, Message: "server address is configured", Evidence: address})

	if err := checkTCP(address, opts.TCPTimeout); err != nil {
		add(CheckResult{Stage: "tcp connectivity", Status: DiagnoseStatusFail, Message: err.Error(), Evidence: address})
		addSkipped(add, "tcp connectivity failed", "login token", "status", "transaction query", "global lock query")
		add(diagnoseDBResult(opts, DiagnoseStatusSkip, "skipped because tcp connectivity failed"))
		finish()
		return snapshot
	}
	add(CheckResult{Stage: "tcp connectivity", Status: DiagnoseStatusPass, Message: "tcp connection succeeded", Evidence: address})

	token, err := GetAuth().GetToken()
	if err != nil {
		add(CheckResult{Stage: "login token", Status: DiagnoseStatusFail, Message: err.Error()})
		addSkipped(add, "login token failed", "status", "transaction query", "global lock query")
		add(diagnoseDBResult(opts, DiagnoseStatusSkip, "skipped because login token failed"))
		finish()
		return snapshot
	}
	add(CheckResult{Stage: "login token", Status: DiagnoseStatusPass, Message: "login token is available", Evidence: maskToken(token)})

	if response, err := QueryStatus(); err != nil {
		add(CheckResult{Stage: "status", Status: DiagnoseStatusFail, Message: err.Error(), Evidence: HealthCheckURL})
		addSkipped(add, "status failed", "transaction query", "global lock query")
		add(diagnoseDBResult(opts, DiagnoseStatusSkip, "skipped because status failed"))
		finish()
		return snapshot
	} else {
		add(CheckResult{Stage: "status", Status: DiagnoseStatusPass, Message: fmt.Sprintf("found %d node(s)", len(response.Data)), Evidence: HealthCheckURL})
	}

	var (
		transactionResponse *GlobalSessionPageResult
		transactionErr      error
		lockResponse        *GlobalLockPageResult
		lockErr             error
	)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		transactionResponse, transactionErr = QueryGlobalSessions(GlobalSessionQuery{PageNum: opts.PageNum, PageSize: opts.PageSize})
	}()
	go func() {
		defer wg.Done()
		lockResponse, lockErr = QueryGlobalLocks(GlobalLockQuery{PageNum: opts.LockPageNum, PageSize: opts.LockPageSize})
	}()
	wg.Wait()

	if transactionErr != nil {
		snapshot.TransactionErr = transactionErr
		add(CheckResult{Stage: "transaction query", Status: DiagnoseStatusFail, Message: transactionErr.Error(), Evidence: GlobalSessionQueryURL})
	} else {
		snapshot.Transactions = transactionResponse
		add(CheckResult{Stage: "transaction query", Status: DiagnoseStatusPass, Message: fmt.Sprintf("found %d global session(s)", len(transactionResponse.Data)), Evidence: GlobalSessionQueryURL})
	}

	if lockErr != nil {
		snapshot.LockErr = lockErr
		add(CheckResult{Stage: "global lock query", Status: DiagnoseStatusFail, Message: lockErr.Error(), Evidence: GlobalLockQueryURL})
	} else {
		snapshot.Locks = lockResponse
		add(CheckResult{Stage: "global lock query", Status: DiagnoseStatusPass, Message: fmt.Sprintf("found %d global lock(s)", len(lockResponse.Data)), Evidence: GlobalLockQueryURL})
	}

	add(diagnoseDBResult(opts, DiagnoseStatusSkip, "database/schema check is optional in this build"))
	finish()
	return snapshot
}

func normalizeDiagnoseOptions(opts DiagnoseOptions) DiagnoseOptions {
	if opts.TCPTimeout <= 0 {
		opts.TCPTimeout = 3 * time.Second
	}
	if opts.PageNum <= 0 {
		opts.PageNum = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 1
	}
	if opts.LockPageNum <= 0 {
		opts.LockPageNum = 1
	}
	if opts.LockPageSize <= 0 {
		opts.LockPageSize = 1
	}
	return opts
}

func addSkipped(add func(CheckResult), reason string, stages ...string) {
	for _, stage := range stages {
		add(CheckResult{Stage: stage, Status: DiagnoseStatusSkip, Message: "skipped because " + reason})
	}
}

func diagnoseDBResult(opts DiagnoseOptions, fallbackStatus string, fallbackMessage string) CheckResult {
	if !opts.CheckDB {
		return CheckResult{Stage: "db/schema", Status: DiagnoseStatusSkip, Message: "database/schema check not requested"}
	}
	if strings.TrimSpace(opts.DBAddress) == "" {
		return CheckResult{Stage: "db/schema", Status: fallbackStatus, Message: fallbackMessage}
	}
	if err := checkTCP(opts.DBAddress, opts.TCPTimeout); err != nil {
		return CheckResult{Stage: "db/schema", Status: DiagnoseStatusFail, Message: err.Error(), Evidence: opts.DBAddress}
	}
	return CheckResult{Stage: "db/schema", Status: DiagnoseStatusPass, Message: "database endpoint is reachable", Evidence: opts.DBAddress}
}

func checkTCP(address string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return err
	}
	return conn.Close()
}

func FormatDiagnoseReport(response *DiagnoseReport, output string) (string, error) {
	output, err := NormalizeOutput(output)
	if err != nil {
		return "", err
	}
	switch output {
	case OutputTable:
		return FormatDiagnoseTable(response.Data), nil
	default:
		return FormatStructuredOutput(response, output)
	}
}

func FormatDiagnoseTable(results []CheckResult) string {
	t := table.NewWriter()
	t.Style().Format.Header = text.FormatDefault
	t.AppendHeader(table.Row{"stage", "status", "message", "detail", "evidence"})
	for _, result := range results {
		t.AppendRow(table.Row{result.Stage, result.Status, result.Message, result.Detail, result.Evidence})
	}
	return t.Render()
}

func maskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 8 {
		return token
	}
	return token[:4] + "..." + token[len(token)-4:]
}

func boolPtr(value bool) *bool {
	return &value
}
