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

package prometheus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/guptarohit/asciigraph"
	"github.com/seata/seata-ctl/action/k8s/utils"
	"github.com/seata/seata-ctl/model"
	"github.com/seata/seata-ctl/tool"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

var MetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Show Prometheus metrics",
	Run: func(_ *cobra.Command, _ []string) {
		if err := showMetrics(); err != nil {
			tool.Logger.Errorf("Failed to show metrics: %v", err)
		}
	},
}

func init() {
	MetricsCmd.PersistentFlags().StringVar(&Target, "target", DefaultPromTarget, "seata prometheus metrics name")
}

// showMetrics executes the metrics collection and chart generation
func showMetrics() error {
	prometheusURL, err := getPrometheusAddress()
	if err != nil {
		return err
	}

	// Query Prometheus for metrics
	result, err := queryPromMetric(prometheusURL, Target)
	if err != nil {
		return fmt.Errorf("query prometheus metrics: %v", err)
	}

	// Generate terminal chart from the queried results
	if err = generateTerminalLineChart(result, Target); err != nil {
		return err
	}
	return nil
}

// getPrometheusAddress fetches Prometheus server address from configuration
func getPrometheusAddress() (string, error) {
	file, err := os.ReadFile(utils.ConfigFileName)
	if err != nil {
		return "", fmt.Errorf("failed to read config.yml: %v", err)
	}

	// Parse the configuration
	var config model.Config
	if err = yaml.Unmarshal(file, &config); err != nil {
		return "", fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	// Extract Prometheus address based on context
	contextName := config.Context.Prometheus
	var contextPath string
	for _, server := range config.Prometheus.Servers {
		if server.Name == contextName {
			contextPath = server.Address
		}
	}
	if contextPath == "" {
		return "", fmt.Errorf("failed to find Prometheus context in config.yml")
	}
	return contextPath, nil
}

// PromResponse defines the structure of a Prometheus query response
type PromResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// queryPromsMetric sends a query to the Prometheus API and returns the response
func queryPromMetric(prometheusURL, query string) (*PromResponse, error) {
	queryURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, url.QueryEscape(query))
	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, fmt.Errorf("error querying Prometheus: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			tool.Logger.Errorf("failed to close response body: %v", err)
			return
		}
	}(resp.Body)

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON response into the PrometheusResponse structure
	var result PromResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	return &result, nil
}

// generateTerminalLineChart generates and prints an ASCII line chart based on the Prometheus response
func generateTerminalLineChart(response *PromResponse, metricName string) error {
	var yValues []float64

	// Iterate over the results and extract the values for the specified metric
	for _, result := range response.Data.Result {
		if name, ok := result.Metric["__name__"]; ok && name == metricName {
			// Parse the metric value
			if valueStr, ok := result.Value[1].(string); ok {
				value, err := strconv.ParseFloat(valueStr, 64)
				if err != nil {
					return fmt.Errorf("error converting value to float: %v", err)
				}
				yValues = append(yValues, value)

			} else {
				return fmt.Errorf("error converting value to float: %v", result.Value[1])
			}
		}
	}

	// Check if any values were found
	if len(yValues) == 0 {
		return fmt.Errorf("no data found for metric: %s", metricName)
	}

	// Generate and display the ASCII line chart
	graph := asciigraph.Plot(yValues, asciigraph.Width(50), asciigraph.Height(10), asciigraph.Caption(metricName))

	tool.Logger.Infof(graph)
	return nil
}
