package logadapter

import (
	"encoding/json"
	"fmt"
	"github.com/seata/seata-ctl/tool"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// QueryLogs sends a request to the /query endpoint and retrieves logs based on the provided filter.
// - filter: A map that holds log filtering parameters such as log level.
// - currency: A struct containing information like the source application and API address.
// - number: The number of logs to fetch (limit).
func (l *Local) QueryLogs(filter map[string]interface{}, currency *Currency, number int) error {
	logLevel, ok := filter["logLevel"].(string)
	if !ok {
		return fmt.Errorf("logLevel is missing or invalid in the filter")
	}

	// Build the query URL using the filter and currency information
	url := fmt.Sprintf("%s%s?application_id=%s&log_level=%s&limit=%d", currency.Address, LocalQueryPath, currency.Source, logLevel, number)

	// Send a GET request to the /query endpoint
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make request to %s: %v", url, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			tool.Logger.Error("failed to close response body: %v", err)
		}
	}(resp.Body)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Check if the response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	// Parse the JSON response body into QueryResponse struct
	var queryResponse QueryResponse
	err = json.Unmarshal(body, &queryResponse)
	if err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Print the logs in a structured format
	for _, value := range queryResponse.Logs {
		res := fmt.Sprintf("[%s]: %s\n", value.Timestamp, value.LogMessage)
		fmt.Println(value.LogMessage)
		if strings.Contains(value.LogMessage, "INFO") {
			tool.Logger.Info(fmt.Sprintf("%v", res))
		}
		if strings.Contains(value.LogMessage, "ERROR") {
			tool.Logger.Error(fmt.Sprintf("%v", res))
		}
		if strings.Contains(value.LogMessage, "WARN") {
			tool.Logger.Warn(fmt.Sprintf("%v", res))
		}
	}
	return nil
}
