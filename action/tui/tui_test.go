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
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/seata/seata-ctl/seata"
)

func TestModelNavigationAndRefreshControls(t *testing.T) {
	m := newModel(false, "", 20, 5)
	if !m.refreshing || m.activePage != pageDiagnose {
		t.Fatalf("initial model = %+v", m)
	}

	updated, cmd := m.Update(key('2'))
	m = updated.(model)
	if cmd != nil || m.activePage != pageTransaction {
		t.Fatalf("transaction navigation = page %d, cmd %v", m.activePage, cmd)
	}

	m.refreshing = false
	updated, cmd = m.Update(key('r'))
	m = updated.(model)
	if !m.refreshing || cmd == nil {
		t.Fatalf("manual refresh = %+v, cmd %v", m, cmd)
	}
	updated, cmd = m.Update(key('r'))
	m = updated.(model)
	if !m.refreshing || cmd != nil {
		t.Fatalf("duplicate refresh = %+v, cmd %v", m, cmd)
	}

	updated, cmd = m.Update(snapshotMsg{})
	m = updated.(model)
	if m.refreshing || cmd == nil {
		t.Fatalf("completed refresh = %+v, cmd %v", m, cmd)
	}

	updated, cmd = m.Update(key('a'))
	m = updated.(model)
	if m.autoRefresh || cmd != nil {
		t.Fatalf("disabled auto refresh = %+v, cmd %v", m, cmd)
	}
	updated, cmd = m.Update(refreshMsg{})
	m = updated.(model)
	if cmd != nil || m.refreshing {
		t.Fatalf("refresh while disabled = %+v, cmd %v", m, cmd)
	}
}

func TestModelKeepsPreviousDataAndShowsWarning(t *testing.T) {
	m := newModel(false, "", 20, 5)
	transactionID := "1001"
	transactions := &seata.GlobalSessionPageResult{
		Data: []seata.GlobalSession{{TransactionID: &transactionID}},
	}
	m.applySnapshot(snapshot{
		Report:       reportWithStage("transaction query", seata.DiagnoseStatusPass, "ok"),
		Transactions: transactions,
	})

	m.applySnapshot(snapshot{
		Report:         reportWithStage("transaction query", seata.DiagnoseStatusFail, "server unavailable"),
		TransactionErr: errors.New("server unavailable"),
	})
	output := m.viewTransaction()
	if !strings.Contains(output, "1001") || !strings.Contains(output, "WARN") {
		t.Fatalf("transaction view = %q, want previous data and warning", output)
	}
}

func TestModelShowsSkippedPageAsUnavailable(t *testing.T) {
	m := newModel(false, "", 20, 5)
	m.applySnapshot(snapshot{
		Report: reportWithStage("global lock query", seata.DiagnoseStatusSkip, "status failed"),
	})
	if output := m.viewLock(); !strings.Contains(output, "unavailable") || !strings.Contains(output, "status failed") {
		t.Fatalf("lock view = %q, want skipped warning", output)
	}
}

func TestModelKeepsLastUpdatedWhenRefreshFails(t *testing.T) {
	m := newModel(false, "", 20, 5)
	old := time.Unix(1710000000, 0)
	m.lastUpdated = old

	updated, cmd := m.Update(snapshotMsg{snapshot: snapshot{Report: reportWithStage("transaction query", seata.DiagnoseStatusFail, "server unavailable")}})
	m = updated.(model)
	if cmd == nil {
		t.Fatal("expected tick command after refresh")
	}
	if !m.lastUpdated.Equal(old) {
		t.Fatalf("lastUpdated = %v, want %v", m.lastUpdated, old)
	}
}

func key(runeValue rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{runeValue}}
}

func reportWithStage(stage, status, message string) *seata.DiagnoseReport {
	success := status == seata.DiagnoseStatusPass
	return &seata.DiagnoseReport{
		Success: &success,
		Data:    []seata.CheckResult{{Stage: stage, Status: status, Message: message}},
	}
}
