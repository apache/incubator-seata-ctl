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

	"github.com/seata/seata-ctl/action"
	"github.com/seata/seata-ctl/action/common"
	"github.com/seata/seata-ctl/seata"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "seata-ctl",
		Short: "seata-ctl is a CLI tool for Seata",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
)

func init() {
	credential := seata.GetAuth()
	rootCmd.PersistentFlags().StringVar(&credential.ServerIP, "ip", "127.0.0.1", "Seata Server IP")
	rootCmd.PersistentFlags().IntVar(&credential.ServerPort, "port", 7091, "Seata Server Admin Port")
	rootCmd.PersistentFlags().StringVar(&credential.Username, "username", "seata", "Username")
	rootCmd.PersistentFlags().StringVar(&credential.Password, "password", "seata", "Password")
	viper.BindPFlag("ip", rootCmd.PersistentFlags().Lookup("ip"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "seata-ctl",
		Short: "seata-ctl is a CLI tool for Seata",
		Run: func(cmd *cobra.Command, args []string) {
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
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" || arg == "version" {
			os.Exit(0)
		}
	}
	//address := seata.GetAuth().GetAddress()
	//err := seata.GetAuth().Login()
	//if err != nil {
	//	fmt.Println("login failed!")
	//	os.Exit(1)
	//}
	var err error

	for {
		//printPrompt(address)
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
