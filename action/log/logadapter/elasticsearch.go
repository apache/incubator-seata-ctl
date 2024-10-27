package logadapter

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/seata/seata-ctl/action/log"
	"github.com/seata/seata-ctl/tool"
	"net/http"
	"strings"
)

type Elasticsearch struct{}

// QueryLogs is a function that queries specific documents
func (e *Elasticsearch) QueryLogs(filter map[string]interface{}, currency *Currency, number int) error {
	client, err := createElasticClient(currency)
	if err != nil {
		return fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	indexName := currency.Source

	// Build the query based on the filter provided
	query, err := BuildQueryFromFilter(filter)
	if err != nil {
		return err
	}

	// Execute the search query
	searchResult, err := client.Search().
		Index(indexName).
		Size(number).
		Query(query).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching documents: %w", err)
	}

	err = processSearchHits(searchResult)
	if err != nil {
		return err
	}
	return nil
}

// createElasticClient configures and creates a new Elasticsearch client
func createElasticClient(currency *Currency) (*elastic.Client, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client, err := elastic.NewClient(
		elastic.SetURL(currency.Address),
		elastic.SetHttpClient(httpClient),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(log.ElasticsearchAuth, currency.Auth),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// processSearchHits handles and formats the search results
func processSearchHits(searchResult *elastic.SearchResult) error {

	if len(searchResult.Hits.Hits) == 0 {
		return fmt.Errorf("no documents found")
	}

	for _, hit := range searchResult.Hits.Hits {
		var doc map[string]interface{}
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			return fmt.Errorf("failed to unmarshal document: %w", err)
		}

		// Pretty print the document content
		for key, value := range doc {
			if key == "log" {
				if strings.Contains(value.(string), "INFO") {
					tool.Logger.Info(fmt.Sprintf("%v", value))
				}
				if strings.Contains(value.(string), "ERROR") {
					tool.Logger.Error(fmt.Sprintf("%v", value))
				}
				if strings.Contains(value.(string), "WARN") {
					tool.Logger.Warn(fmt.Sprintf("%v", value))
				}
			}
		}
	}
	return nil
}

// BuildQueryFromFilter constructs a BoolQuery from the filter parameters
func BuildQueryFromFilter(filter map[string]interface{}) (*elastic.BoolQuery, error) {
	query := elastic.NewBoolQuery()

	for key, value := range filter {
		switch key {
		case "logLevel":
			if v, ok := value.(string); ok {
				query.Should(elastic.NewTermQuery("logLevel", v))
			} else {
				return nil, fmt.Errorf("invalid type for logLevel, expected string")
			}
		case "module":
			if v, ok := value.(string); ok {
				query.Should(elastic.NewTermQuery("module", v))
			} else {
				return nil, fmt.Errorf("invalid type for module, expected string")
			}
		case "xid":
			if v, ok := value.(string); ok {
				query.Should(elastic.NewTermQuery("xid", v))
			} else {
				return nil, fmt.Errorf("invalid type for xid, expected string")
			}
		case "branchId":
			if v, ok := value.(string); ok {
				query.Should(elastic.NewTermQuery("branchId", v))
			} else {
				return nil, fmt.Errorf("invalid type for branchId, expected string")
			}
		case "resourceId":
			if v, ok := value.(string); ok {
				query.Should(elastic.NewTermQuery("resourceId", v))
			} else {
				return nil, fmt.Errorf("invalid type for resourceId, expected string")
			}
		case "message":
			if v, ok := value.(string); ok {
				query.Should(elastic.NewMatchQuery("message", v))
			} else {
				return nil, fmt.Errorf("invalid type for message, expected string")
			}
		default:
			return nil, fmt.Errorf("unknown field: %s", key)
		}
	}
	return query, nil
}
