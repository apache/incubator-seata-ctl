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

package cmd

import (
	"fmt"
	"os"

	"github.com/seata/seata-ctl/action/login"
	"github.com/seata/seata-ctl/tool"

	"github.com/seata/seata-ctl/action"
	"github.com/seata/seata-ctl/action/common"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "seata-ctl",
		Short: "seata-ctl is a CLI tool for Seata",
		Run: func(_ *cobra.Command, _ []string) {
			// Do Stuff Here
		},
	}
)

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "seata-ctl",
		Short: "seata-ctl is a CLI tool for Seata",
		Run: func(_ *cobra.Command, _ []string) {
			//
		},
	})
	rootCmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
		HiddenDefaultCmd:    true,
	}

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tool.InitLogger()

	var address = ""

	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" || arg == "version" {
			os.Exit(0)
		}
	}
	var err error
	for {
		if login.Address != "" {
			printPrompt(address)
		}
		err = common.ReadArgs(os.Stdin)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err = action.Execute(); err != nil {
			fmt.Println(err)
			os.Args = []string{}
		}
	}
}

func printPrompt(address string) {
	fmt.Printf("%s > ", address)
}
