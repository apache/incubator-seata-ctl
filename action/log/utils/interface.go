package utils

type LogQuery interface {
	// QueryLogs 根据查询条件获取日志
	QueryLogs(filter map[string]interface{}) ([]string, error)
}

type Currency struct {
	Address string `json:"address"`
	Source  string `json:"source"`
}
