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

package login

import (
	"fmt"
	"os"

	"github.com/seata/seata-ctl/seata"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Address string

// LoginCmd 定义 login 子命令
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Seata server",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Attempting to login...")
		Address = seata.GetAuth().GetAddress()
		err := seata.GetAuth().Login()
		if err != nil {
			fmt.Println("Login failed!")
			os.Exit(1)
		}
		fmt.Printf("Login successful to address: %s\n", Address)
		printPrompt(Address)
	},
}

func init() {
	credential := seata.GetAuth()
	LoginCmd.PersistentFlags().StringVar(&credential.ServerIP, "ip", "127.0.0.1", "Seata Server IP")
	LoginCmd.PersistentFlags().IntVar(&credential.ServerPort, "port", 7091, "Seata Server Admin Port")
	LoginCmd.PersistentFlags().StringVar(&credential.Username, "username", "seata", "Username")
	LoginCmd.PersistentFlags().StringVar(&credential.Password, "password", "seata", "Password")
	err := viper.BindPFlag("ip", LoginCmd.PersistentFlags().Lookup("ip"))
	if err != nil {
		return
	}
	err = viper.BindPFlag("port", LoginCmd.PersistentFlags().Lookup("port"))
	if err != nil {
		return
	}
	err = viper.BindPFlag("username", LoginCmd.PersistentFlags().Lookup("username"))
	if err != nil {
		return
	}
	err = viper.BindPFlag("password", LoginCmd.PersistentFlags().Lookup("password"))
	if err != nil {
		return
	}
}

func printPrompt(address string) {
	fmt.Printf("%s > ", address)
}
