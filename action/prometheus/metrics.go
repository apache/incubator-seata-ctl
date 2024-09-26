package prometheus

import (
	"encoding/json"
	"fmt"
	"github.com/seata/seata-ctl/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var MetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "show prometheus metrics",
	Run: func(cmd *cobra.Command, args []string) {
		err := metrics()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var Target string

func init() {
	MetricsCmd.PersistentFlags().StringVar(&Target, "target", "seata_transaction_summary", "Namespace name")
}

func metrics() error {
	prometheusURL, err := getPrometheusAddress()
	if err != nil {
		return err
	}
	result, err := queryPrometheusMetric(prometheusURL, Target)
	if err != nil {
		log.Fatalf("Error querying Prometheus: %v", err)
	}
	printPrometheusResponse(result)
	return nil
}

// QueryResult 用于解析 Prometheus 查询结果的结构体
type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func getPrometheusAddress() (string, error) {
	file, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Failed to read config.yml: %v", err)
	}
	var config model.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	contextName := config.Context.Prometheus
	var contextPath string
	for _, server := range config.Prometheus.Servers {
		if server.Name == contextName {
			contextPath = server.Address
		}
	}
	if contextPath == "" {
		log.Fatalf("Failed to find context in config.yml")
		return "", err
	}
	return contextPath, err
}

//func generateView(result *PrometheusResponse) error {
//	app := tview.NewApplication()
//	table := tview.NewTable().
//		SetBorders(true).
//		SetFixed(1, 0)
//	headers := []string{"Metric", "Timestamp", "Value"}
//	for i, header := range headers {
//		table.SetCell(0, i, tview.NewTableCell(header).
//			SetTextColor(tview.TrueColor).
//			SetAlign(tview.AlignCenter).
//			SetSelectable(false))
//	}
//	for rowIndex, res := range result.Data.Result {
//		value := res.Value
//		if len(value) == 2 {
//			timestamp, ok := value[0].(float64)
//			if !ok {
//				log.Println("Invalid timestamp format")
//				continue
//			}
//			val, ok := value[1].(string)
//			if !ok {
//				log.Println("Invalid value format")
//				continue
//			}
//			timeStamp := time.Unix(int64(timestamp), 0).UTC()
//			table.SetCell(rowIndex+1, 0, tview.NewTableCell(Target).
//				SetAlign(tview.AlignLeft))
//			table.SetCell(rowIndex+1, 1, tview.NewTableCell(timeStamp.Format(time.RFC3339)).
//				SetAlign(tview.AlignCenter))
//			table.SetCell(rowIndex+1, 2, tview.NewTableCell(val).
//				SetAlign(tview.AlignRight))
//		}
//	}
//	if err := app.SetRoot(table, true).Run(); err != nil {
//		log.Fatalf("Error starting tview application: %v", err)
//	}
//	return nil
//}

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

func queryPrometheusMetric(prometheusURL, query string) (*PrometheusResponse, error) {
	queryURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, url.QueryEscape(query))
	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, fmt.Errorf("error querying Prometheus: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	var result PrometheusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	return &result, nil
}

func printPrometheusResponse(response *PrometheusResponse) {
	if response == nil {
		fmt.Println("The Prometheus response is nil.")
		return
	}

	fmt.Println("Prometheus Response Status:", response.Status)
	if response.Data.ResultType != "" {
		fmt.Println("Result Type:", response.Data.ResultType)
	}

	if len(response.Data.Result) > 0 {
		fmt.Println("Results:")
		for i, result := range response.Data.Result {
			fmt.Printf("  Result %d:\n", i+1)
			if len(result.Metric) > 0 {
				fmt.Printf("    Metric: %+v\n", result.Metric)
			}
			if len(result.Value) > 0 {
				fmt.Printf("    Value: %+v\n", result.Value)
			}
		}
	} else {
		fmt.Println("  No results found.")
	}
}
