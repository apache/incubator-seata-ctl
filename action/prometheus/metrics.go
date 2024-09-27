package prometheus

import (
	"encoding/json"
	"fmt"
	"github.com/guptarohit/asciigraph"
	"github.com/seata/seata-ctl/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var MetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Show Prometheus metrics",
	Run: func(cmd *cobra.Command, args []string) {
		if err := showMetrics(); err != nil {
			fmt.Println(err)
		}
	},
}

var Target string

func init() {
	MetricsCmd.PersistentFlags().StringVar(&Target, "target", "seata_transaction_summary", "Namespace name")
}

// showMetrics executes the metrics collection and chart generation
func showMetrics() error {
	prometheusURL, err := getPrometheusAddress()
	if err != nil {
		return err
	}

	// Query Prometheus for metrics
	result, err := queryPrometheusMetric(prometheusURL, Target)
	if err != nil {
		log.Fatalf("Error querying Prometheus: %v", err)
	}

	// Generate terminal chart from the queried results
	if err = generateTerminalLineChart(result, Target); err != nil {
		return err
	}
	return nil
}

// getPrometheusAddress fetches Prometheus server address from configuration
func getPrometheusAddress() (string, error) {
	file, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Failed to read config.yml: %v", err)
	}

	// Parse the configuration
	var config model.Config
	if err = yaml.Unmarshal(file, &config); err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
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
		log.Fatalf("Failed to find Prometheus context in config.yml")
		return "", err
	}
	return contextPath, nil
}

// PrometheusResponse defines the structure of a Prometheus query response
type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// queryPrometheusMetric sends a query to the Prometheus API and returns the response
func queryPrometheusMetric(prometheusURL, query string) (*PrometheusResponse, error) {
	queryURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, url.QueryEscape(query))
	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, fmt.Errorf("error querying Prometheus: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON response into the PrometheusResponse structure
	var result PrometheusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	return &result, nil
}

// generateTerminalLineChart generates and prints an ASCII line chart based on the Prometheus response
func generateTerminalLineChart(response *PrometheusResponse, metricName string) error {
	var yValues []float64

	// Iterate over the results and extract the values for the specified metric
	for _, result := range response.Data.Result {
		if name, ok := result.Metric["__name__"]; ok && name == metricName {
			// Parse the metric value
			if valueStr, ok := result.Value[1].(string); ok {
				value, err := strconv.ParseFloat(valueStr, 64)
				if err != nil {
					fmt.Println("Value: Invalid number format")
				} else {
					yValues = append(yValues, value)
				}
			} else {
				fmt.Println("Value: Invalid value format")
			}
		}
	}

	// Check if any values were found
	if len(yValues) == 0 {
		return fmt.Errorf("no data found for metric: %s", metricName)
	}

	// Generate and display the ASCII line chart
	graph := asciigraph.Plot(yValues, asciigraph.Width(50), asciigraph.Height(10), asciigraph.Caption(metricName))
	fmt.Println(graph)
	return nil
}
