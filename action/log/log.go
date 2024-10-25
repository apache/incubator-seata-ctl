package log

import (
	"fmt"
	"github.com/seata/seata-ctl/action/log/logadapter"
	"github.com/seata/seata-ctl/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const (
	ElasticSearchType = "ElasticSearch"
	DefaultNumber     = 10
	DefaultLogLevel   = ""
)

const (
	LokiType = "Loki"
)

const (
	LocalType = "Local"
)

var LogCmd = &cobra.Command{
	Use:   "log",
	Short: "get seata log",
	Run: func(cmd *cobra.Command, args []string) {
		err := getLog()
		if err != nil {
			fmt.Println(err)
		}
	},
}

// ElasticSearch
var (
	LogLevel   string
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

func init() {
	LogCmd.PersistentFlags().StringVar(&LogLevel, "level", DefaultLogLevel, "seata log level")
	LogCmd.PersistentFlags().StringVar(&Module, "module", "", "seata module")
	LogCmd.PersistentFlags().StringVar(&XID, "xid", "", "seata expression")
	LogCmd.PersistentFlags().StringVar(&BranchID, "banchID", "", "seata branchId")
	LogCmd.PersistentFlags().StringVar(&ResourceID, "resourceID", "", "seata resourceID")
	LogCmd.PersistentFlags().StringVar(&Message, "message", "", "seata message")
	LogCmd.PersistentFlags().IntVar(&Number, "number", DefaultNumber, "seata number")
	LogCmd.PersistentFlags().StringVar(&Label, "label", "", "seata label")
	LogCmd.PersistentFlags().StringVar(&Start, "start", "", "seata start")
	LogCmd.PersistentFlags().StringVar(&End, "end", "", "seata end")
}

func getLog() error {
	context, currency, err := getContext()
	if err != nil {
		return err
	}
	logType := context.Types

	var client logadapter.LogQuery
	var filter = make(map[string]interface{})

	switch logType {
	case ElasticSearchType:
		{
			client = &logadapter.Elasticsearch{}
			filter = buildElasticSearchFilter()
		}
	case LokiType:
		{
			client = &logadapter.Loki{}
			filter = buildLokiFilter()
		}
	case LocalType:
		{
			client = &logadapter.Local{}
			filter = buildLocalFilter()
		}
	}

	if client == nil {
		return fmt.Errorf("can not get client")
	}

	err = client.QueryLogs(filter, currency, Number)
	if err != nil {
		return err
	}

	//reset var
	//ResetAllVariables()

	return nil
}

func getContext() (*model.Cluster, *logadapter.Currency, error) {
	file, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Failed to read config.yml: %v", err)
	}
	var config model.Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	contextName := config.Context.Log
	for _, cluster := range config.Log.Clusters {
		if cluster.Name == contextName {
			currency := logadapter.Currency{
				Address: cluster.Address,
				Source:  cluster.Source,
				Auth:    cluster.Auth,
			}
			return &cluster, &currency, nil
		}
	}
	return nil, nil, fmt.Errorf("failed to find context in config.yml")
}

func buildElasticSearchFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	if LogLevel != "" {
		filter["logLevel"] = LogLevel
	}
	if Module != "" {
		filter["module"] = Module
	}
	if XID != "" {
		filter["XID"] = XID
	}
	if BranchID != "" {
		filter["BranchID"] = BranchID
	}
	if ResourceID != "" {
		filter["ResourceID"] = ResourceID
	}
	if Message != "" {
		filter["message"] = Message
	}
	return filter
}

func buildLokiFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	filter["query"] = Label
	if Start != "" {
		filter["start"] = Start
	}
	if End != "" {
		filter["end"] = End
	}
	return filter
}

func buildLocalFilter() map[string]interface{} {
	filter := make(map[string]interface{})
	if LogLevel != "" {
		filter["logLevel"] = LogLevel
	}
	return filter
}

// ResetAllVariables resets all global variables to their zero values
func ResetAllVariables() {
	// Reset ElasticSearch-related variables
	LogLevel = ""
	Module = ""
	XID = ""
	BranchID = ""
	ResourceID = ""
	Message = ""
	Number = 0

	// Reset Loki-related variables
	Label = ""
	Start = ""
	End = ""
}
