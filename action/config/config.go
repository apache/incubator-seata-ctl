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
	Short: "Set config path",
	Run: func(cmd *cobra.Command, args []string) {
		err := createYMLFile(Path)
		if err != nil {
			println("Error creating config:", err.Error())
			log.Fatal(err)
		}
		println("Config created successfully!")
	},
}

func init() {
	ConfigCmd.PersistentFlags().StringVar(&Path, "path", "/", "Set config path")
}

func createSampleConfig() model.Config {
	return model.Config{
		Kubernetes: model.Kubernetes{
			Cluster: []model.KubernetesCluster{
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
		Context: model.Context{
			Kubernetes: "",
			Prometheus: "",
			Log:        "",
		},
	}
}

// Create a YAML file
func createYMLFile(path string) error {
	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	// Check if the file already exists
	ymlFilePath := path + "config.yml"
	if _, err := os.Stat(ymlFilePath); err == nil {
		return fmt.Errorf("file already exists: %s", path)
	}

	// Create a sample config object
	config := createSampleConfig()

	// Marshal the config object into YAML format
	data, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("Failed to marshal config to YAML: %v", err)
	}

	// Write the YAML data to a file
	file, err := os.Create("config.yml")
	if err != nil {
		log.Fatalf("Failed to create YAML config: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("Failed to close config file: %v", err)
		}
	}()

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("Failed to write to YAML config: %v", err)
	}
	fmt.Println("YAML config created successfully")

	// Update the path variable
	Path = path
	return nil
}
