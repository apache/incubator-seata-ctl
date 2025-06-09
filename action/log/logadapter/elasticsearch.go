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

package logadapter

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	elastic "github.com/olivere/elastic/v7"
	"github.com/seata/seata-ctl/tool"
)

// QueryLogs is a function that queries specific documents
func (e *Elasticsearch) QueryLogs(filter map[string]interface{}, currency *Currency, number int) error {
	client, err := createElasticClient(currency)
	if err != nil {
		return fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	indexName := currency.Source

	indexFields, err := getEsIndexList(currency)
	if err != nil {
		return err
	}
	query, err := buildQuery(filter, indexFields)
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

	err = processSearchHits(searchResult, currency)
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
		elastic.SetBasicAuth(currency.Username, currency.Password),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// createEsDefaultClient configures and creates a new Elasticsearch client
func createEsDefaultClient(currency *Currency) (*elasticsearch.Client, error) {
	// Configure the Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{
			currency.Address,
		},
		Username: currency.Username,
		Password: currency.Password,
		// Skip certificate verification if using a self-signed certificate
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Create the client instance
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the client: %s", err)
	}
	return es, nil
}

// processSearchHits handles and formats the search results
func processSearchHits(searchResult *elastic.SearchResult, currency *Currency) error {
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
			if key == currency.Index {
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

// getFieldNames recursively extracts field names under the "fields" key
func getFieldNames(properties map[string]interface{}, prefix string) []string {
	fieldNames := []string{}

	for fieldName, fieldValue := range properties {
		// Generate the full path for the current field
		fullName := fieldName
		if prefix != "" {
			fullName = prefix + "." + fieldName
		}

		// Check if the field contains a "fields" node
		if fieldMap, ok := fieldValue.(map[string]interface{}); ok {
			if fields, ok := fieldMap["fields"].(map[string]interface{}); ok {
				// If there is a "fields" node, iterate through its fields and add to the result
				for subField := range fields {
					fieldNames = append(fieldNames, fullName+"."+subField)
				}
			}

			// If the field contains nested "properties", recursively parse subfields
			if nestedProperties, ok := fieldMap["properties"].(map[string]interface{}); ok {
				fieldNames = append(fieldNames, getFieldNames(nestedProperties, fullName)...)
			}
		}
	}

	return fieldNames
}

// extractFields extracts all field names from a nested map structure
func extractFields(data map[string]interface{}) []string {
	var allFields []string

	// Iterate through each index to get its field names
	for _, indexData := range data {
		if indexMap, ok := indexData.(map[string]interface{}); ok {
			if mappings, ok := indexMap["mappings"].(map[string]interface{}); ok {
				if properties, ok := mappings["properties"].(map[string]interface{}); ok {
					// Get all field names under "fields" and merge into the result
					allFields = append(allFields, getFieldNames(properties, "")...)
				}
			}
		}
	}

	return allFields
}

// ParseJobString parses the input string and returns a map
func ParseJobString(input string) (map[string]string, error) {
	// Remove curly braces
	input = strings.Trim(input, "{}")

	// Split by ','
	parts := strings.Split(input, ",")
	kvMap := make(map[string]string)

	for _, part := range parts {
		// Split by '=' to get key-value pairs
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid key=value pair: %s", part)
		}
		kvMap[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	return kvMap, nil
}

// Contains checks if a string exists in a slice of strings
func Contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// getEsIndexList retrieves field names from the specified Elasticsearch index
func getEsIndexList(currency *Currency) ([]string, error) {
	es, err := createEsDefaultClient(currency)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	// Build the request to get the mappings
	req := esapi.IndicesGetMappingRequest{
		Index: []string{currency.Source}, // Specify the index name
	}

	// Execute the request
	res, err := req.Do(context.Background(), es)
	if err != nil {
		return nil, fmt.Errorf("error getting mapping: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			tool.Logger.Error("error closing body")
		}
	}(res.Body)

	// Check if the response is successful
	if res.IsError() {
		return nil, fmt.Errorf("error response: %s", res.String())
	}

	// Read and parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to make request to %s", err)
	}

	// Call method to extract field names
	indexFields := extractFields(result)
	indexFields = removeKeywordSuffix(indexFields)
	return indexFields, nil
}

// removeKeywordSuffix removes ".keyword" suffix from each string in the slice
func removeKeywordSuffix(input []string) []string {
	var result []string
	for _, str := range input {
		// Check if the string ends with ".keyword"
		str = strings.TrimSuffix(str, ".keyword")
		result = append(result, str) // Add the processed string to the result slice
	}
	return result
}

// buildQuery constructs a BoolQuery based on the provided filter and index fields
func buildQuery(filter map[string]interface{}, indexFields []string) (*elastic.BoolQuery, error) {
	query := elastic.NewBoolQuery()
	if filter["query"].(string) != "{}" {
		indexMap, err := ParseJobString(filter["query"].(string))
		if err != nil {
			return query, err
		}
		for k, v := range indexMap {
			if Contains(indexFields, k) {
				query.Should(elastic.NewTermQuery(k, v))
			} else {
				return query, fmt.Errorf("invalid index key: %s", k)
			}
		}
	}
	return query, nil
}
