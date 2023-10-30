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

package try

import (
	"github.com/seata/seata-ctl/action/common"
	"github.com/seata/seata-ctl/seata"
	"github.com/spf13/cobra"
)

var (
	rollbackXID string
)

func init() {
	RollbackCmd.SetUsageTemplate(common.GetUsageTmpl("rollback"))
	RollbackCmd.SetHelpTemplate(common.GetHelpTmpl())
	RollbackCmd.PersistentFlags().StringVar(&rollbackXID, "xid", "", "rollback a txn with xid")
}

var RollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "rollback a txn",
	Run: func(cmd *cobra.Command, args []string) {
		seata.RollbackTxn(rollbackXID)
		rollbackXID = ""
	},
}
