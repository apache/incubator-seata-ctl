package config

import (
	"fmt"
	"github.com/seata/seata-ctl/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var Path = "/"

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "config config path",
	Run: func(cmd *cobra.Command, args []string) {
		err := createYMLFile(Path)
		if err != nil {
			println("create config error:", err.Error())
			log.Fatal(err)
		}
		println("create config success!")
	},
}

func init() {
	ConfigCmd.PersistentFlags().StringVar(&Path, "path", "/", "config config path")
}

func createSampleConfig() model.Config {
	return model.Config{
		Kubernetes: model.Kubernetes{
			Clusters: []model.KubernetesCluster{
				{
					Name:           "",
					KubeConfigPath: "",
					YmlPath:        "",
				},
			},
		},
		Prometheus: model.Prometheus{
			Servers: []model.Server{
				{
					Name:    "",
					Address: "",
					Auth:    "",
				},
			},
		},
		Log: model.Log{
			Clusters: []model.Cluster{
				{
					Name: "",
					Collection: model.Collection{
						Enable: true,
						Local:  "",
					},
					Analysis: model.Analysis{
						Enable: true,
						Local:  "",
					},
					Display: model.Display{
						DisplayType: "",
						Path:        "",
						Local:       "",
					},
				},
			},
		},
	}
}

// 创建yml文件
func createYMLFile(path string) error {
	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	// 检查文件是否存在
	ymlFilePath := path + "/config.yml"
	if _, err := os.Stat(ymlFilePath); err == nil {
		return fmt.Errorf("file already exists: %s", path)
	}

	// 创建示例配置对象
	config := createSampleConfig()

	// 将配置对象编码为 YAML 格式
	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("Failed to marshal config to YAML: %v", err)
	}

	// 将 YAML 数据写入文件
	file, err := os.Create("config.yml")
	if err != nil {
		log.Fatalf("Failed to create YAML config: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("Failed to close config: %v", err)
		}
	}()

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("Failed to write to YAML config: %v", err)
	}
	fmt.Println("YAML config created successfully")

	//修改path的位置
	Path = path

	return nil
}
