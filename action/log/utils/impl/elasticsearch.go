package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/seata/seata-ctl/action/log/utils"
	"log"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

type LogEntry struct {
	Timestamp    time.Time `json:"timestamp,omitempty"`     // 使用 time.Time 来存储 ISO 8601 格式的时间戳
	LogLevel     string    `json:"log_level,omitempty"`     // 日志级别，例如 "INFO"
	Thread       string    `json:"thread,omitempty"`        // 线程名称
	LoggerClass  string    `json:"logger_class,omitempty"`  // 产生日志的类名
	LoggerMethod string    `json:"logger_method,omitempty"` // 产生日志的方法名
	LineNumber   int       `json:"line_number,omitempty"`   // 代码的行号
	Message      string    `json:"message,omitempty"`       // 实际的日志消息
	XID          string    `json:"XID,omitempty"`           // 事务的唯一标识
}

type Elasticsearch struct {
}

type ElasticsearchParams struct {
	LogLevel      string `json:"log_level"`
	TransactionId string `json:"transaction_id"`
}

type Match struct {
	Field map[string]string `json:"match"`
}

type BoolQuery struct {
	Must []Match `json:"must"`
}

type Query struct {
	Bool BoolQuery `json:"bool"`
}

type SearchQuery struct {
	Query Query `json:"query"`
}

func newSearchQuery(logLevel, transactionID string) (string, error) {
	must := []Match{
		{Field: map[string]string{"log_level": logLevel}},
		{Field: map[string]string{"transaction_id": transactionID}},
	}
	boolQuery := BoolQuery{
		Must: must,
	}
	query := Query{
		Bool: boolQuery,
	}
	searchQuery := SearchQuery{
		Query: query,
	}
	queryJSON, err := json.Marshal(searchQuery)
	if err != nil {
		return "", err
	}
	return string(queryJSON), nil
}

func (e *Elasticsearch) QueryLogs(filter map[string]interface{}) ([]string, error) {
	currency, ok := filter["currency"].(utils.Currency)
	if !ok {
		return nil, fmt.Errorf("error: failed to assert currency type")
	}
	params, ok := filter["Elasticsearch"].(ElasticsearchParams)
	if !ok {
		return nil, fmt.Errorf("error: failed to assert elasticsearchParams type")
	}
	cfg := elasticsearch.Config{
		Addresses: []string{currency.Address},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	query, err := newSearchQuery(params.LogLevel, params.TransactionId)
	if err != nil {
		return nil, err
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(currency.Source),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	var logStrings []string
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var entry LogEntry
		source := hit.(map[string]interface{})["_source"]
		sourceData, _ := json.Marshal(source)
		if err := json.Unmarshal(sourceData, &entry); err != nil {
			log.Fatalf("Error unmarshalling log: %s", err)
		}
		logString := fmt.Sprintf("Timestamp: %s, LogLevel: %s, Thread: %s, LoggerClass: %s, LoggerMethod: %s, LineNumber: %d, Message: %s, XID: %s",
			entry.Timestamp.Format(time.RFC3339),
			entry.LogLevel,
			entry.Thread,
			entry.LoggerClass,
			entry.LoggerMethod,
			entry.LineNumber,
			entry.Message,
			entry.XID,
		)
		logStrings = append(logStrings, logString)
	}
	return logStrings, nil
}
