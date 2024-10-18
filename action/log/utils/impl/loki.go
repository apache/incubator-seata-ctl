package impl

import (
	"encoding/json"
	"fmt"
	"github.com/seata/seata-ctl/action/log/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	LokiAddressPath = "/loki/api/v1/query_range?"
)

type Loki struct{}

// LokiQueryResult represents the structure of the response from Loki
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

// QueryLogs queries logs from Loki based on the filter and settings provided
func (l *Loki) QueryLogs(filter map[string]interface{}, currency *utils.Currency, number int) error {

	// Prepare the query URL with time range
	params := url.Values{}
	params.Set("query", filter["query"].(string))
	params.Set("limit", strconv.Itoa(filter["number"].(int)))

	// Set start time if provided
	if value, ok := filter["start"]; ok {
		res, err := parseToTimestamp(value.(string))
		if err != nil {
			return err
		}
		params.Set("start", fmt.Sprintf("%d", res))
	}
	// Set end time if provided
	if value, ok := filter["end"]; ok {
		res, err := parseToTimestamp(value.(string))
		if err != nil {
			return err
		}
		params.Set("end", fmt.Sprintf("%d", res))
	}
	queryURL := currency.Address + LokiAddressPath + params.Encode()

	// Send GET request to Loki
	resp, err := http.Get(queryURL)
	if err != nil {
		log.Fatalf("Error sending request to Loki: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response from Loki: %v", err)
	}

	// Parse the JSON response
	var result LokiQueryResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Error parsing Loki response: %v", err)
	}

	// Process and print logs
	if result.Status == "success" {
		for _, stream := range result.Data.Result {
			for _, entry := range stream.Values {
				// Extract timestamp and log message
				timestampStr := entry[0]
				logLine := entry[1]

				// Convert nanosecond timestamp to int64
				timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
				if err != nil {
					fmt.Printf("Failed to convert timestamp: %v\n", err)
					continue
				}

				// Convert Unix nanosecond timestamp to readable format
				readableTime := time.Unix(0, timestampInt).Format("2006-01-02 15:04:05")

				// Print the readable timestamp and log message
				fmt.Printf("Timestamp: %s, Log: %s\n", readableTime, logLine)
			}
		}
	} else {
		fmt.Printf("Query failed, status: %s\n", result.Status)
	}
	return nil
}

// parseToTimestamp parses a time string to a Unix nanosecond timestamp
func parseToTimestamp(timeStr string) (int64, error) {
	// Define the time layout to match the input format
	const timeLayout = "2006-01-02-15:04:05"

	// Load the specified time zone (UTC in this case)
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return 0, fmt.Errorf("failed to load timezone: %v", err)
	}

	// Parse the time string using the specified layout and timezone
	t, err := time.ParseInLocation(timeLayout, timeStr, loc)
	if err != nil {
		return 0, fmt.Errorf("failed to parse time: %v", err)
	}

	// Convert to Unix nanosecond timestamp
	timestamp := t.UnixNano()

	return timestamp, nil
}
