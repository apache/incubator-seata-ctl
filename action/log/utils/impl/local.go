package impl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/seata/seata-ctl/action/log/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Local struct{}

// 日志响应结构
type LogsResponse struct {
	Logs map[string]string `json:"logs"` // 日志内容以节点名为键，日志内容为值
}

func (l *Local) QueryLogs(filter map[string]interface{}) ([]string, error) {

	currency, ok := filter["currency"].(utils.Currency)
	if !ok {
		return nil, fmt.Errorf("error: failed to assert currency type")
	}
	params, ok := filter["Local"].(LokiParams)
	if !ok {
		return nil, fmt.Errorf("error: failed to assert lokiParams type")
	}
	url := fmt.Sprintf("%s/logs/%s", currency.Address, params.Expression)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var logsResponse LogsResponse
	err = json.Unmarshal(body, &logsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	var allLogs []string
	for node, logs := range logsResponse.Logs {
		scanner := bufio.NewScanner(strings.NewReader(logs))
		for scanner.Scan() {
			line := scanner.Text()
			formattedLine := fmt.Sprintf("Node: %s - %s", node, line)
			allLogs = append(allLogs, formattedLine)
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Failed to scan logs for node %s: %v", node, err)
		}
	}
	return nil, nil
}
