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

package action

import (
	"github.com/seata/seata-ctl/action/common"
	"github.com/seata/seata-ctl/action/get"
	"github.com/seata/seata-ctl/action/reload"
	se "github.com/seata/seata-ctl/action/set"
	del "github.com/seata/seata-ctl/action/try"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(quitCmd,
		get.GetCmd, del.TryCmd,
		reload.ReloadCmd, se.SetCmd)
	rootCmd.SetHelpTemplate(common.GetHelpTmplWithOnlyAvailableCmd())
	rootCmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
		HiddenDefaultCmd:    true,
	}
}

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() error {
	return rootCmd.Execute()
}
