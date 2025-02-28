package k8s

const (
	DefaultCRName          = "example-seataserver"
	DefaultServerImage     = "apache/seata-server:latest"
	DefaultNamespace       = "default"
	DefaultDeployName      = "seata-k8s-controller-manager"
	DefaultControllerImage = "apache/seata-controller:latest"
	DefaultReplicas        = 1

	Label          = "cr_name"
	CreateCrdPath  = "/apis/apiextensions.k8s.io/v1/customresourcedefinitions"
	FilePath       = "seata.yaml"
	CRDname        = "seataservers.operator.seata.apache.org"
	ServiceName    = "seata-server-cluster"
	RequestStorage = "1Gi"
	LimitStorage   = "1Gi"
)

var (
	Name            string
	Replicas        int32
	Namespace       string
	Image           string
	ControllerImage string
	DeployName      string
)
