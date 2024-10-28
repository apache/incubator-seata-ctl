package logadapter

import "time"

type Elasticsearch struct{}

type Loki struct{}

type Local struct{}

type Currency struct {
	Address  string `json:"address"`
	Source   string `json:"source"`
	Username string `json:"username"`
	Password string `json:"password"`
	Index    string `json:"index"`
}

type SeataLog struct {
	Timestamp  time.Time // Timestamp
	LogLevel   string    // Log level, e.g., INFO, ERROR, etc.
	Module     string    // Log module, e.g., RMHandler, UndoLogManager, etc.
	XID        string    // Global transaction ID
	BranchID   string    // Branch transaction ID
	ResourceID string    // Resource ID
	Message    string    // Log content
}

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

// LogData represents a single log entry
type LogData struct {
	Timestamp  string `json:"timestamp"`
	LogLevel   string `json:"log_level"`
	LogMessage string `json:"log_message"`
}

// QueryResponse holds the response structure from the /query API
type QueryResponse struct {
	ApplicationID string    `json:"application_id"`
	LogLevel      string    `json:"log_level"`
	Logs          []LogData `json:"logs"`
}
