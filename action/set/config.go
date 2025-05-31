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

package se

import (
	"github.com/seata/seata-ctl/action/common"
	"github.com/seata/seata-ctl/seata"
	"github.com/spf13/cobra"
)

var (
	kvData          string
	setRegistry     bool
	setConfigCenter bool
)

func init() {
	ConfigCmd.PersistentFlags().StringVar(&kvData, "data", "{}", "Configuration map")
	ConfigCmd.PersistentFlags().BoolVar(&setRegistry, "registry", false, "If set registry conf")
	ConfigCmd.PersistentFlags().BoolVar(&setConfigCenter, "config-center", false, "If set configuration center conf")
}

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Set the configuration",
	Run: func(_ *cobra.Command, _ []string) {
		data, err := common.ParseDictArg(kvData)
		if err != nil {
			common.Log("", err)
		}
		configType := seata.NormalConfig
		if setRegistry {
			configType = seata.RegistryConf
		} else if setConfigCenter {
			configType = seata.ConfigCenterConf
		}
		common.Log(seata.SetConfiguration(data, configType))
		kvData = "{}"
		setRegistry = false
		setConfigCenter = false
	},
}
