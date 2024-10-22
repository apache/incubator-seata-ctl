package logadapter

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"net/http"
)

const (
	ElasticsearchAuth = "elastic"
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

	fmt.Printf("Found %d hits.\n", searchResult.TotalHits())
	processSearchHits(searchResult)
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
		elastic.SetBasicAuth(ElasticsearchAuth, currency.Auth),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// processSearchHits handles and formats the search results
func processSearchHits(searchResult *elastic.SearchResult) []string {
	var result []string

	for _, hit := range searchResult.Hits.Hits {
		var doc map[string]interface{}
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			log.Printf("Error parsing document: %s", err)
			continue
		}

		// Pretty print the document content
		fmt.Println("Document Source:")
		for key, value := range doc {
			log.Printf("  %s: %v", key, value)
		}
		log.Println("----------------------------")

		result = append(result, hit.Id) // Store document IDs or other relevant data
	}
	return result
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
				return nil, errors.New("invalid type for logLevel, expected string")
			}
		case "module":
			if v, ok := value.(string); ok {
				query.Should(elastic.NewTermQuery("module", v))
			} else {
				return nil, errors.New("invalid type for module, expected string")
			}
		case "xid":
			if v, ok := value.(string); ok {
				query.Should(elastic.NewTermQuery("xid", v))
			} else {
				return nil, errors.New("invalid type for xid, expected string")
			}
		case "branchId":
			if v, ok := value.(string); ok {
				query.Should(elastic.NewTermQuery("branchId", v))
			} else {
				return nil, errors.New("invalid type for branchId, expected string")
			}
		case "resourceId":
			if v, ok := value.(string); ok {
				query.Should(elastic.NewTermQuery("resourceId", v))
			} else {
				return nil, errors.New("invalid type for resourceId, expected string")
			}
		case "message":
			if v, ok := value.(string); ok {
				query.Should(elastic.NewMatchQuery("message", v))
			} else {
				return nil, errors.New("invalid type for message, expected string")
			}
		default:
			log.Printf("Unknown field: %s\n", key)
		}
	}
	return query, nil
}
