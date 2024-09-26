package utils

import "fmt"

func CreateRequest(createCrdPath string, filePath string) (string, error) {

	kubeConfigPath, err := GetKubeConfigPath()
	if err != nil {
		return "", fmt.Errorf("failed to get kubeconfig path: %v", err)
	}
	kubeConfig, err := LoadKubeConfig(kubeConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to load kubeconfig file: %v", err)
	}
	contextInfo, err := GetContextInfo(kubeConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get context info: %v", err)
	}

	filePath, err = ConvertAndSaveYamlToJSON(filePath)

	if err != nil {
		return "", fmt.Errorf("failed to save yaml: %v", err)
	}

	res, err := sendPostRequest(contextInfo, createCrdPath, filePath)
	if err != nil {
		return res, fmt.Errorf("failed to send post request: %v", err)
	}
	return res, nil
}
