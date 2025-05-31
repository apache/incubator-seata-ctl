package config

import (
	"fmt"
	"log"
	"os"

	"github.com/seata/seata-ctl/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var Path string

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Set config path",
	Run: func(_ *cobra.Command, _ []string) {
		err := createYMLFile(Path)
		if err != nil {
			println("Error creating config:", err.Error())
			log.Fatal(err)
		}
	},
}

func init() {
	ConfigCmd.PersistentFlags().StringVar(&Path, "path", "", "Set config path")
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
					Name:     "",
					Types:    "",
					Address:  "",
					Source:   "",
					Username: "",
					Password: "",
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
	//if _, err := os.Stat(path); os.IsNotExist(err) {
	//	return fmt.Errorf("path does not exist: %s", path)
	//}

	// Check if the file already exists
	ymlFilePath := "config.yml"
	if _, err := os.Stat(ymlFilePath); err == nil {
		fmt.Println("Config file already exists!")
		return nil
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
	fmt.Println("Config created successfully!")
	// Update the path variable
	Path = path
	return nil
}
