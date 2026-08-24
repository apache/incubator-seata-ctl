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

package transaction

import (
	"fmt"
	"strings"

	"github.com/seata/seata-ctl/action/common"
	"github.com/seata/seata-ctl/seata"
	"github.com/spf13/cobra"
)

var ListCmd = newListCommand()

func newListCommand() *cobra.Command {
	var (
		xid             string
		applicationID   string
		status          int
		transactionName string
		vgroup          string
		withBranch      bool
		pageNum         int
		pageSize        int
		timeStart       int64
		timeEnd         int64
		output          string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List global transactions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			defer common.ResetLocalFlags(cmd)

			output = strings.ToLower(output)
			if output != seata.OutputTable && output != seata.OutputJSON && output != seata.OutputYAML {
				return fmt.Errorf("unsupported output format %q", output)
			}

			query := seata.GlobalSessionQuery{
				XID:             xid,
				ApplicationID:   applicationID,
				TransactionName: transactionName,
				Vgroup:          vgroup,
				PageNum:         pageNum,
				PageSize:        pageSize,
			}
			if cmd.Flags().Changed("status") {
				query.Status = &status
			}
			if cmd.Flags().Changed("with-branch") {
				query.WithBranch = &withBranch
			}
			if cmd.Flags().Changed("time-start") {
				query.TimeStart = &timeStart
			}
			if cmd.Flags().Changed("time-end") {
				query.TimeEnd = &timeEnd
			}

			response, err := seata.QueryGlobalSessions(query)
			if err != nil {
				return err
			}
			result, err := seata.FormatGlobalSessionPage(response, output)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
	common.ResetLocalFlagsOnParseErrorAndHelp(cmd)
	cmd.SetUsageTemplate(common.GetUsageTmpl("transaction list"))
	cmd.SetHelpTemplate(common.GetHelpTmpl())
	cmd.Flags().StringVar(&xid, "xid", "", "Filter by transaction XID")
	cmd.Flags().StringVar(&applicationID, "application-id", "", "Filter by application ID")
	cmd.Flags().IntVar(&status, "status", 0, "Filter by global transaction status code")
	cmd.Flags().StringVar(&transactionName, "transaction-name", "", "Filter by transaction name")
	cmd.Flags().StringVar(&vgroup, "vgroup", "", "Filter by transaction service group")
	cmd.Flags().BoolVar(&withBranch, "with-branch", false, "Include branch sessions")
	cmd.Flags().IntVar(&pageNum, "page-num", 1, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "Page size")
	cmd.Flags().Int64Var(&timeStart, "time-start", 0, "Filter by begin time start, in milliseconds")
	cmd.Flags().Int64Var(&timeEnd, "time-end", 0, "Filter by begin time end, in milliseconds")
	cmd.Flags().StringVar(&output, "output", seata.OutputTable, "Output format: table, json, yaml")
	return cmd
}
