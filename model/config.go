package model

type Config struct {
	Kubernetes Kubernetes `yaml:"kubernetes"`
	Prometheus Prometheus `yaml:"prometheus"`
	Log        Log        `yaml:"log"`
}

type Kubernetes struct {
	Cluster KubernetesCluster `yaml:"clusters"`
}

type Prometheus struct {
	Servers []Server `yaml:"servers,omitempty"`
}

type Log struct {
	Clusters []Cluster `yaml:"clusters"`
}

type KubernetesCluster struct {
	Name           string `yaml:"name"`
	KubeConfigPath string `yaml:"kubeconfigpath"`
	YmlPath        string `yaml:"ymlpath"`
}

type Server struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Auth    string `yaml:"auth"`
}

type Cluster struct {
	Name       string     `yaml:"name"`
	Collection Collection `yaml:"collection"`
	Analysis   Analysis   `yaml:"analysis"`
	Display    Display    `yaml:"display"`
}

type Collection struct {
	Enable bool   `yaml:"enable"`
	Local  string `yaml:"local"`
}

type Analysis struct {
	Enable bool   `yaml:"enable"`
	Local  string `yaml:"local"`
}

type Display struct {
	DisplayType string `yaml:"displayType"`
	Path        string `yaml:"path"`
	Local       string `yaml:"local"`
}
