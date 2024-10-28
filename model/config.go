package model

type Config struct {
	Kubernetes Kubernetes `yaml:"kubernetes"`
	Prometheus Prometheus `yaml:"prometheus"`
	Log        Log        `yaml:"log"`
	Context    Context    `yaml:"context"`
}

type Kubernetes struct {
	Cluster []KubernetesCluster `yaml:"clusters"`
}

type Prometheus struct {
	Servers []Server `yaml:"servers,omitempty"`
}

type Log struct {
	Clusters []Cluster `yaml:"clusters"`
}

// Context Select the appropriate configuration based on the Context field
type Context struct {
	Kubernetes string `yaml:"kubernetes"`
	Prometheus string `yaml:"prometheus"`
	Log        string `yaml:"log"`
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
	Name     string `yaml:"name"`
	Types    string `yaml:"types"`
	Address  string `yaml:"address"`
	Source   string `yaml:"source"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Index    string `yaml:"index"`
}
