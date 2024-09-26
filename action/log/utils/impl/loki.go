package impl

import (
	"encoding/json"
	"fmt"
	"github.com/seata/seata-ctl/action/log/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Loki struct {
}

type LokiParams struct {
	Expression string `json:"expression"`
}

type LokiQueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func ConvertLokiResultToStrings(result LokiQueryResult) ([]string, error) {
	var logs []string
	for _, res := range result.Data.Result {
		for _, value := range res.Values {
			timestampNano, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse timestamp: %w", err)
			}
			timestamp := time.Unix(0, timestampNano).Format(time.RFC3339)
			logLine := value[1]
			logs = append(logs, fmt.Sprintf("[%s] %s", timestamp, logLine))
		}
	}
	return logs, nil
}

func (l *Loki) QueryLogs(filter map[string]interface{}) ([]string, error) {
	currency, ok := filter["currency"].(utils.Currency)
	if !ok {
		return nil, fmt.Errorf("error: failed to assert currency type")
	}
	params, ok := filter["Loki"].(LokiParams)
	if !ok {
		return nil, fmt.Errorf("error: failed to assert lokiParams type")
	}

	// 构建 Loki 查询 URL
	queryParams := url.Values{}
	queryParams.Add("query", currency.Source+params.Expression)
	queryURL := fmt.Sprintf("%s/loki/api/v1/query_range?%s", currency.Address, queryParams.Encode())

	// 发送 HTTP 请求
	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 解析 JSON
	var result LokiQueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	//转换格式
	res, err := ConvertLokiResultToStrings(result)
	if err != nil {
		return nil, err
	}
	return res, nil
}
