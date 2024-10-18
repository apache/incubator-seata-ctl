package utils

import "time"

type LogQuery interface {
	// QueryLogs retrieves logs based on the filter criteria
	QueryLogs(filter map[string]interface{}, currency *Currency, number int) error
}

type Currency struct {
	Address string `json:"address"`
	Source  string `json:"source"`
	Auth    string `json:"auth"`
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
