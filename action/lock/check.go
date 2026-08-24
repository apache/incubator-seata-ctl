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

var CheckCmd = newCheckCommand()

func newCheckCommand() *cobra.Command {
	var (
		xid      string
		branchID string
		output   string
	)

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check whether a global lock exists",
		RunE: func(cmd *cobra.Command, _ []string) error {
			defer common.ResetLocalFlags(cmd)

			output = strings.ToLower(output)
			if _, err := seata.NormalizeOutput(output); err != nil {
				return err
			}

			response, err := seata.CheckGlobalLock(xid, branchID)
			if err != nil {
				return err
			}
			result, err := seata.FormatGlobalLockCheck(response, output)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}
	common.ResetLocalFlagsOnParseErrorAndHelp(cmd)
	cmd.SetUsageTemplate(common.GetUsageTmpl("lock check"))
	cmd.SetHelpTemplate(common.GetHelpTmpl())
	cmd.Flags().StringVar(&xid, "xid", "", "Transaction XID")
	cmd.Flags().StringVar(&branchID, "branch-id", "", "Branch ID")
	cmd.Flags().StringVar(&output, "output", seata.OutputTable, "Output format: table, json, yaml")
	return cmd
}
