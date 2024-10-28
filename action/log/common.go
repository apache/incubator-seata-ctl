package log

const (
	ElasticSearchType = "ElasticSearch"
	LokiType          = "Loki"
	LocalType         = "Local"

	DefaultNumber        = 10
	DefaultLogLevel      = ""
	DefaultLocalLogLevel = "-"
)

var (
	Label  string
	Number int
	Start  string
	End    string
	Level  string
)
