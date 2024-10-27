package log

const (
	ElasticSearchType = "ElasticSearch"
	LokiType          = "Loki"
	LocalType         = "Local"

	DefaultNumber        = 10
	DefaultLogLevel      = ""
	DefaultLocalLogLevel = "-"

	LokiAddressPath = "/loki/api/v1/query_range?"
	TimeLayout      = "2006-01-02-15:04:05"

	LocalQueryPath = "/query"

	ElasticsearchAuth = "elastic"
)

// ElasticSearch
var (
	Level      string
	Module     string
	XID        string
	BranchID   string
	ResourceID string
	Message    string
	Number     int
)

// Loki
var (
	Label string
	Start string
	End   string
)
