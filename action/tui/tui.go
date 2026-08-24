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

package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/seata/seata-ctl/action/common"
	"github.com/seata/seata-ctl/seata"
	"github.com/spf13/cobra"
)

func init() {
	TuiCmd.SetUsageTemplate(common.GetUsageTmpl("tui"))
	TuiCmd.SetHelpTemplate(common.GetHelpTmpl())
	TuiCmd.Flags().Duration("refresh", 5*time.Second, "Auto refresh interval")
	TuiCmd.Flags().Int("page-size", 20, "Transaction and lock page size")
	TuiCmd.Flags().Bool("check-db", false, "Check database connectivity")
	TuiCmd.Flags().String("db-address", "", "Database address in host:port form")
}

var TuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Open the diagnostic TUI",
	RunE: func(cmd *cobra.Command, _ []string) error {
		refreshEvery, _ := cmd.Flags().GetDuration("refresh")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		checkDB, _ := cmd.Flags().GetBool("check-db")
		dbAddress, _ := cmd.Flags().GetString("db-address")
		return tea.NewProgram(newModel(checkDB, dbAddress, pageSize, refreshEvery), tea.WithAltScreen()).Start()
	},
}

type page int

const (
	pageDiagnose page = iota
	pageTransaction
	pageLock
)

type snapshot = seata.DiagnoseSnapshot

type snapshotMsg struct{ snapshot snapshot }
type refreshMsg struct{}

type model struct {
	activePage   page
	autoRefresh  bool
	refreshEvery time.Duration
	pageSize     int
	checkDB      bool
	dbAddress    string
	snapshot     snapshot
	refreshing   bool
	lastUpdated  time.Time
}

func newModel(checkDB bool, dbAddress string, pageSize int, refreshEvery time.Duration) model {
	if pageSize <= 0 {
		pageSize = 20
	}
	if refreshEvery <= 0 {
		refreshEvery = 5 * time.Second
	}
	return model{
		activePage:   pageDiagnose,
		autoRefresh:  true,
		refreshEvery: refreshEvery,
		pageSize:     pageSize,
		checkDB:      checkDB,
		dbAddress:    dbAddress,
		refreshing:   true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.refreshCmd(), m.tickCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.activePage = pageDiagnose
		case "2":
			m.activePage = pageTransaction
		case "3":
			m.activePage = pageLock
		case "tab":
			m.activePage = (m.activePage + 1) % 3
		case "shift+tab":
			m.activePage = (m.activePage + 2) % 3
		case "a":
			m.autoRefresh = !m.autoRefresh
			return m, m.tickCmd()
		case "r":
			if m.refreshing {
				return m, nil
			}
			m.refreshing = true
			return m, m.refreshCmd()
		}
	case tea.WindowSizeMsg:
		return m, nil
	case refreshMsg:
		if !m.autoRefresh || m.refreshing {
			return m, nil
		}
		m.refreshing = true
		return m, m.refreshCmd()
	case snapshotMsg:
		m.refreshing = false
		m.applySnapshot(msg.snapshot)
		if msg.snapshot.Report != nil && msg.snapshot.Report.Success != nil && *msg.snapshot.Report.Success {
			m.lastUpdated = time.Now()
		}
		return m, m.tickCmd()
	}
	return m, nil
}

func (m model) View() string {
	var body string
	switch m.activePage {
	case pageDiagnose:
		body = m.viewDiagnose()
	case pageTransaction:
		body = m.viewTransaction()
	case pageLock:
		body = m.viewLock()
	}

	state := "manual"
	if m.autoRefresh {
		state = "auto"
	}
	updated := "loading"
	if !m.lastUpdated.IsZero() {
		updated = m.lastUpdated.Format("15:04:05")
	}
	return fmt.Sprintf("seata-ctl tui | [1] diagnose [2] transaction [3] lock | refresh=%s | %s | %s\n\n%s\n", m.refreshEvery, state, updated, body)
}

func (m model) refreshCmd() tea.Cmd {
	return func() tea.Msg {
		return snapshotMsg{snapshot: collectSnapshot(m.checkDB, m.dbAddress, m.pageSize)}
	}
}

func (m model) tickCmd() tea.Cmd {
	if !m.autoRefresh || m.refreshEvery <= 0 {
		return nil
	}
	return tea.Tick(m.refreshEvery, func(time.Time) tea.Msg { return refreshMsg{} })
}

func (m *model) applySnapshot(next snapshot) {
	if next.Report != nil {
		m.snapshot.Report = next.Report
	}
	if next.Transactions != nil {
		m.snapshot.Transactions = next.Transactions
	}
	if next.TransactionErr != nil {
		m.snapshot.TransactionErr = next.TransactionErr
	} else {
		m.snapshot.TransactionErr = nil
	}
	if next.Locks != nil {
		m.snapshot.Locks = next.Locks
	}
	if next.LockErr != nil {
		m.snapshot.LockErr = next.LockErr
	} else {
		m.snapshot.LockErr = nil
	}
}

func (m model) viewDiagnose() string {
	if m.snapshot.Report == nil {
		return "loading diagnostics..."
	}
	return strings.TrimSpace(renderDiagnosePage(m.snapshot.Report))
}

func (m model) viewTransaction() string {
	issue := diagnoseStageIssue(m.snapshot.Report, "transaction query")
	if m.snapshot.Transactions == nil {
		if m.snapshot.TransactionErr != nil {
			return "transaction query failed: " + m.snapshot.TransactionErr.Error()
		}
		if issue != "" {
			return "transaction query unavailable: " + issue
		}
		return "loading transactions..."
	}
	output := strings.TrimSpace(seata.FormatGlobalSessionTable(m.snapshot.Transactions.Data))
	if issue != "" {
		return output + "\n\nWARN: transaction query " + issue
	}
	return output
}

func (m model) viewLock() string {
	issue := diagnoseStageIssue(m.snapshot.Report, "global lock query")
	if m.snapshot.Locks == nil {
		if m.snapshot.LockErr != nil {
			return "lock query failed: " + m.snapshot.LockErr.Error()
		}
		if issue != "" {
			return "lock query unavailable: " + issue
		}
		return "loading locks..."
	}
	output := strings.TrimSpace(seata.FormatGlobalLockTable(m.snapshot.Locks.Data))
	if issue != "" {
		return output + "\n\nWARN: lock query " + issue
	}
	return output
}

func collectSnapshot(checkDB bool, dbAddress string, pageSize int) snapshot {
	return *seata.CollectDiagnoseSnapshot(seata.DiagnoseOptions{
		CheckDB:      checkDB,
		DBAddress:    dbAddress,
		TCPTimeout:   3 * time.Second,
		PageNum:      1,
		PageSize:     pageSize,
		LockPageNum:  1,
		LockPageSize: pageSize,
	})
}

func renderDiagnosePage(report *seata.DiagnoseReport) string {
	output, err := seata.FormatDiagnoseReport(report, seata.OutputTable)
	if err != nil {
		return err.Error()
	}
	if report.Success != nil && !*report.Success {
		return output + "\n\nWARN: diagnostics report has failures"
	}
	return output
}

func diagnoseStageIssue(report *seata.DiagnoseReport, stage string) string {
	if report == nil {
		return ""
	}
	for _, result := range report.Data {
		if result.Stage == stage && result.Status != seata.DiagnoseStatusPass {
			return result.Status + ": " + result.Message
		}
	}
	return ""
}
