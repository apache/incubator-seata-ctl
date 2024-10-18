package impl

type Local struct{}

// 日志响应结构
type LogsResponse struct {
	Logs map[string]string `json:"logs"` // 日志内容以节点名为键，日志内容为值
}

func (l *Local) QueryLogs(filter map[string]interface{}) ([]string, error) {
	return nil, nil
}
