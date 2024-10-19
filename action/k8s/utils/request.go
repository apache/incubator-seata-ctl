package utils

import "fmt"

func CreateRequest(createCrdPath string, filePath string) (string, error) {

	kubeConfigPath, err := GetKubeConfigPath()
	if err != nil {
		return "", fmt.Errorf("Failed to get kubeconfig path: %v", err)
	}
	kubeConfig, err := LoadKubeConfig(kubeConfigPath)
	if err != nil {
		return "", fmt.Errorf("Failed to load kubeconfig file: %v", err)
	}
	contextInfo, err := GetContextInfo(kubeConfig)
	if err != nil {
		return "", fmt.Errorf("Failed to get context info: %v", err)
	}

	filePath, err = ConvertAndSaveYamlToJSON(filePath)

	if err != nil {
		return "", fmt.Errorf("Failed to save yaml: %v", err)
	}

	res, err := sendPostRequest(contextInfo, createCrdPath, filePath)
	if err != nil {
		return res, fmt.Errorf("Failed to send post request: %v", err)
	}
	return res, nil
}
