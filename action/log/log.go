package log

import (
	"errors"
	"fmt"
	"github.com/seata/seata-ctl/action/log/utils"
	"github.com/seata/seata-ctl/action/log/utils/impl"
	"github.com/seata/seata-ctl/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
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

var Level string
var TransactionId string
var Expression string

func init() {
	LogCmd.PersistentFlags().StringVar(&Level, "level", "ERROR", "seata log level")
	LogCmd.PersistentFlags().StringVar(&TransactionId, "Id", "", "seata transaction id")
	LogCmd.PersistentFlags().StringVar(&Expression, "Expression", "", "seata expression")
}

func getLog() error {
	context, currency, err := getContext()
	if err != nil {
		return err
	}

	param := make(map[string]interface{})
	param["currency"] = currency

	client, param, err := getClientAndParams(context, param)
	if err != nil {
		return err
	}

	res, err := client.QueryLogs(param)
	if err != nil {
		return err
	}

	err = showLogInfo(res)
	if err != nil {
		return err
	}

	return nil
}

func getContext() (*model.Cluster, *utils.Currency, error) {
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
			currency := utils.Currency{
				Address: cluster.Address,
				Source:  cluster.Source,
			}
			return &cluster, &currency, nil
		}
	}
	return nil, nil, fmt.Errorf("failed to find context in config.yml")
}

func getClientAndParams(cluster *model.Cluster, param map[string]interface{}) (utils.LogQuery, map[string]interface{}, error) {
	logType := cluster.Types
	switch logType {
	case "ElasticSearch":
		{
			param["ElasticSearch"] = impl.ElasticsearchParams{
				LogLevel:      logType,
				TransactionId: TransactionId,
			}
			return &impl.Elasticsearch{}, param, nil
		}
	case "Loki":
		{
			param["Loki"] = impl.LokiParams{
				Expression: Expression,
			}
			return &impl.Loki{}, nil, nil
		}
	case "Local":
		{
			param["Local"] = impl.LokiParams{
				Expression: Expression,
			}
			return &impl.Local{}, nil, nil
		}
	default:
		{
			return nil, nil, fmt.Errorf("unknown log type: %s", logType)
		}
	}
}

func showLogInfo(logs []string) error {
	if len(logs) == 0 {
		return errors.New("no logs to display")
	}
	fmt.Println("=========== LOG ENTRIES ===========")
	for i, logEntry := range logs {
		fmt.Printf("Log #%d:\n", i+1)
		fmt.Println("------------------------------")
		fmt.Printf("%s\n", logEntry)
		fmt.Println("------------------------------")
	}
	fmt.Println("=========== END OF LOGS ===========")
	return nil
}
