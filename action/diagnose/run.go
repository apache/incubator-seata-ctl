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

package diagnose

import (
	"fmt"
	"strings"
	"time"

	"github.com/seata/seata-ctl/action/common"
	"github.com/seata/seata-ctl/seata"
	"github.com/spf13/cobra"
)

var RunCmd = newRunCommand()

func newRunCommand() *cobra.Command {
	var (
		output        string
		checkDB       bool
		dbAddress     string
		tcpTimeout    time.Duration
		pageNum       int
		pageSize      int
		lockPageNum   int
		lockPageSize  int
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run diagnostics against the current Seata server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			defer common.ResetLocalFlags(cmd)

			output = strings.ToLower(output)
			if _, err := seata.NormalizeOutput(output); err != nil {
				return err
			}

			report, err := seata.RunDiagnostics(seata.DiagnoseOptions{
				CheckDB:      checkDB,
				DBAddress:    dbAddress,
				TCPTimeout:   tcpTimeout,
				PageNum:      pageNum,
				PageSize:     pageSize,
				LockPageNum:  lockPageNum,
				LockPageSize: lockPageSize,
			})
			if err != nil {
				return err
			}

			result, err := seata.FormatDiagnoseReport(report, output)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
	common.ResetLocalFlagsOnParseErrorAndHelp(cmd)
	cmd.SetUsageTemplate(common.GetUsageTmpl("diagnose run"))
	cmd.SetHelpTemplate(common.GetHelpTmpl())
	cmd.Flags().BoolVar(&checkDB, "check-db", false, "Check database connectivity")
	cmd.Flags().StringVar(&dbAddress, "db-address", "", "Database address in host:port form")
	cmd.Flags().DurationVar(&tcpTimeout, "tcp-timeout", 3*time.Second, "TCP timeout")
	cmd.Flags().IntVar(&pageNum, "page-num", 1, "Transaction page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 1, "Transaction page size")
	cmd.Flags().IntVar(&lockPageNum, "lock-page-num", 1, "Lock page number")
	cmd.Flags().IntVar(&lockPageSize, "lock-page-size", 1, "Lock page size")
	cmd.Flags().StringVar(&output, "output", seata.OutputTable, "Output format: table, json, yaml")
	return cmd
}
