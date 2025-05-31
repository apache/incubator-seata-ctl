package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/seata/seata-ctl/tool"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type ContextType string

const (
	JSON ContextType = "application/json"
)

// ContextInfo is a structure to hold client certificate, key, CA certificate, API server URL, and content type.
type ContextInfo struct {
	ClientCert  string
	ClientKey   string
	CACert      string
	APIServer   string
	ContentType ContextType
}

// LoadKubeConfig loads the kubeconfig from the provided path and filename
func LoadKubeConfig(kubeconfigFullPath string) (*api.Config, error) {
	// Read the kubeconfig file
	kubeconfigBytes, err := os.ReadFile(kubeconfigFullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig file: %v", err)
	}
	// Convert the kubeconfig file content into the Config struct
	config, err := clientcmd.Load(kubeconfigBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %v", err)
	}
	return config, nil
}

// GetContextInfo extracts the client certificate, key, CA certificate, API server URL, and content type
// from the provided *api.Config for the current context.
func GetContextInfo(config *api.Config) (*ContextInfo, error) {
	// Get the current context name
	currentContext := config.CurrentContext
	if currentContext == "" {
		return nil, fmt.Errorf("no current context is set")
	}
	// Get the current context object
	context, ok := config.Contexts[currentContext]
	if !ok {
		return nil, fmt.Errorf("context %s not found", currentContext)
	}
	// Get the cluster object
	cluster, ok := config.Clusters[context.Cluster]
	if !ok {
		return nil, fmt.Errorf("cluster %s not found", context.Cluster)
	}
	// Get the user object
	authInfo, ok := config.AuthInfos[context.AuthInfo]
	if !ok {
		return nil, fmt.Errorf("auth info %s not found", context.AuthInfo)
	}
	// Extract the certificate and API server information
	clientCert := authInfo.ClientCertificate
	clientKey := authInfo.ClientKey
	caCert := cluster.CertificateAuthority
	apiServer := cluster.Server
	// Content type for API request
	contentType := JSON
	// Check if all required fields are present
	if clientCert == "" || clientKey == "" || caCert == "" || apiServer == "" {
		missingFields := []string{}
		if clientCert == "" {
			missingFields = append(missingFields, "client certificate")
		}
		if clientKey == "" {
			missingFields = append(missingFields, "client key")
		}
		if caCert == "" {
			missingFields = append(missingFields, "CA certificate")
		}
		if apiServer == "" {
			missingFields = append(missingFields, "API server")
		}
		return nil, fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}
	// Return the ContextInfo struct
	return &ContextInfo{
		ClientCert:  clientCert,
		ClientKey:   clientKey,
		CACert:      caCert,
		APIServer:   apiServer,
		ContentType: contentType,
	}, nil
}

// sendPostRequest sends a POST request to the Kubernetes API server to create a custom resource definition (CRD)
func sendPostRequest(context *ContextInfo, createCrdPath string, filePath string) (string, error) {
	certFile := context.ClientCert
	keyFile := context.ClientKey
	caCertFile := context.CACert
	url := context.APIServer + createCrdPath
	contentType := context.ContentType

	// Load client certificate and key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return "", fmt.Errorf("failed to load certificate and key: %v", err)
	}

	// Read CA certificate
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return "", fmt.Errorf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return "", fmt.Errorf("failed to append CA certificate")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Read data from file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read data file: %v", err)
	}

	// Create a POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set content type
	req.Header.Set("Content-Type", string(contentType))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			tool.Logger.Error("failed to close response body: %v", err)
		}
	}(resp.Body)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Return appropriate response based on status code
	if resp.StatusCode == http.StatusCreated {
		tool.Logger.Infof("Create seata crd success")
		return "", nil
	}
	if resp.StatusCode == http.StatusConflict {
		return "", fmt.Errorf("seata crd already exists")
	}
	return "error: " + string(body), err
}
