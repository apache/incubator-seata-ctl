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

package log

import (
	"fmt"
	"os"

	"github.com/seata/seata-ctl/action/log/logadapter"
	"github.com/seata/seata-ctl/model"
	"github.com/seata/seata-ctl/tool"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "get seata log",
	Run: func(_ *cobra.Command, _ []string) {
		err := getLog()
		if err != nil {
			tool.Logger.Errorf("get log error: %s", err)
		}
	},
}

func init() {
	LogCmd.PersistentFlags().StringVar(&Level, "level", DefaultLogLevel, "seata log level")
	LogCmd.PersistentFlags().IntVar(&Number, "number", DefaultNumber, "seata number")
	LogCmd.PersistentFlags().StringVar(&Label, "label", "{}", "seata label")
	LogCmd.PersistentFlags().StringVar(&Start, "start", "", "seata start")
	LogCmd.PersistentFlags().StringVar(&End, "end", "", "seata end")
}

func getLog() error {
	context, currency, err := getContext()
	if err != nil {
		return err
	}
	logType := context.Types

	var client logadapter.LogQuery
	var filter = make(map[string]interface{})

	switch logType {
	case ElasticSearchType:
		{
			client = &logadapter.Elasticsearch{}
			filter = buildElasticSearchFilter()
		}
	case LokiType:
		{
			client = &logadapter.Loki{}
			filter = buildLokiFilter()
		}
	case LocalType:
		{
			client = &logadapter.Local{}
			filter = buildLocalFilter()
		}
	}

	if client == nil {
		return fmt.Errorf("can not get client")
	}

	err = client.QueryLogs(filter, currency, Number)
	if err != nil {
		return err
	}

	return nil
}

func getContext() (*model.Cluster, *logadapter.Currency, error) {
	file, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, nil, fmt.Errorf("read config.yml error: %s", err)
	}
	var config model.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, nil, fmt.Errorf("unmarshal config.yml error: %s", err)
	}
	contextName := config.Context.Log
	for _, cluster := range config.Log.Clusters {
		if cluster.Name == contextName {
			currency := logadapter.Currency{
				Address:  cluster.Address,
				Source:   cluster.Source,
				Username: cluster.Username,
				Password: cluster.Password,
				Index:    cluster.Index,
			}
			return &cluster, &currency, nil
		}
	}
	return nil, nil, fmt.Errorf("failed to find context in config.yml")
}

func buildElasticSearchFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter["query"] = Label
	return filter
}

func buildLokiFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter["query"] = Label
	if Start != "" {
		filter["start"] = Start
	}
	if End != "" {
		filter["end"] = End
	}
	return filter
}

func buildLocalFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	if Level != "" {
		filter["logLevel"] = Level
	} else {
		//reset the constants value
		filter["logLevel"] = DefaultLocalLogLevel
	}
	return filter
}
