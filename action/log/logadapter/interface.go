package logadapter

type LogQuery interface {
	// QueryLogs retrieves logs based on the filter criteria
	QueryLogs(filter map[string]interface{}, currency *Currency, number int) error
}
