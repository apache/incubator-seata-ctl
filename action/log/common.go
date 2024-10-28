package log

const (
	ElasticSearchType = "ElasticSearch"
	LokiType          = "Loki"
	LocalType         = "Local"

	DefaultNumber        = 10
	DefaultLogLevel      = ""
	DefaultLocalLogLevel = "-"
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
