package prometheus

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/url"
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

var Name string
var Target string
var Region string

func init() {
	MetricsCmd.PersistentFlags().StringVar(&Name, "name", "prometheus", "Prometheus name")
	MetricsCmd.PersistentFlags().StringVar(&Target, "target", "default", "Namespace name")
	MetricsCmd.PersistentFlags().StringVar(&Region, "region", "default", "Namespace name")
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

func metrics() error {
	// 构建查询的 URL
	prometheusURL := "http://localhost:9090"
	u, _ := url.ParseRequestURI(prometheusURL)
	u.Path = "/api/v1/query"
	// 添加查询参数
	q := u.Query()
	query := "seata.tm.commit"
	q.Set("query", query)
	u.RawQuery = q.Encode()

	// 发送 HTTP GET 请求
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解析响应体为 QueryResult 结构
	var result QueryResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	return err
	//// 检查 Prometheus 查询状态
	//if result.Status != "success" {
	//	return nil, fmt.Errorf("query failed: %s", result.Status)
	//}
	//
	//return &result, nil
}
