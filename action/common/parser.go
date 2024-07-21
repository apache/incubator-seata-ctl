package common

import (
	"fmt"
	"github.com/seata/seata-ctl/model"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// ReadYMLFile 读取并解析 YAML 文件
func ReadYMLFile(path string) (model.Config, error) {
	var config model.Config

	// 检查文件是否存在
	//ymlFilePath := path + "config.yml"
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		return config, fmt.Errorf("file does not exist: %s", path)
	}

	// 读取文件内容
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read file: %v", err)
	}

	// 解析 YAML 内容
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	return config, nil
}
