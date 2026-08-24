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

package lock

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
		xid           string
		tableName     string
		transactionID string
		branchID      string
		pk            string
		resourceID    string
		pageNum       int
		pageSize      int
		timeStart     int64
		timeEnd       int64
		output        string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List global locks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			defer common.ResetLocalFlags(cmd)

			output = strings.ToLower(output)
			if _, err := seata.NormalizeOutput(output); err != nil {
				return err
			}

			query := seata.GlobalLockQuery{
				XID:           xid,
				TableName:     tableName,
				TransactionID: transactionID,
				BranchID:      branchID,
				PK:            pk,
				ResourceID:    resourceID,
				PageNum:       pageNum,
				PageSize:      pageSize,
			}
			if cmd.Flags().Changed("time-start") {
				query.TimeStart = &timeStart
			}
			if cmd.Flags().Changed("time-end") {
				query.TimeEnd = &timeEnd
			}

			response, err := seata.QueryGlobalLocks(query)
			if err != nil {
				return err
			}
			result, err := seata.FormatGlobalLockPage(response, output)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
	common.ResetLocalFlagsOnParseErrorAndHelp(cmd)
	cmd.SetUsageTemplate(common.GetUsageTmpl("lock list"))
	cmd.SetHelpTemplate(common.GetHelpTmpl())
	cmd.Flags().StringVar(&xid, "xid", "", "Filter by transaction XID")
	cmd.Flags().StringVar(&tableName, "table-name", "", "Filter by table name")
	cmd.Flags().StringVar(&transactionID, "transaction-id", "", "Filter by transaction ID")
	cmd.Flags().StringVar(&branchID, "branch-id", "", "Filter by branch ID")
	cmd.Flags().StringVar(&pk, "pk", "", "Filter by primary key")
	cmd.Flags().StringVar(&resourceID, "resource-id", "", "Filter by resource ID")
	cmd.Flags().IntVar(&pageNum, "page-num", 1, "Page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 20, "Page size")
	cmd.Flags().Int64Var(&timeStart, "time-start", 0, "Filter by lock time start, in milliseconds")
	cmd.Flags().Int64Var(&timeEnd, "time-end", 0, "Filter by lock time end, in milliseconds")
	cmd.Flags().StringVar(&output, "output", seata.OutputTable, "Output format: table, json, yaml")
	return cmd
}
