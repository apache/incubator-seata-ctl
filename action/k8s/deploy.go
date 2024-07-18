package k8s

import (
	"github.com/seata/seata-ctl/action/common"
	"github.com/seata/seata-ctl/action/config"
	"github.com/spf13/cobra"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy seata in kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		println("create config success!")
	},
}

var version string

func init() {
	DeployCmd.PersistentFlags().StringVar(&version, "version", "lastest", "crd version")
}

func deploy(version string) error {
	//获取config
	Config, err := common.ReadYMLFile(config.Path)
	if err != nil {
		return err
	}
	//从config中获取kubeconfig
	kubeconfigpath := Config.Kubernetes.Cluster.KubeConfigPath

	//创建客户端
	client, err := common.GetClient(kubeconfigpath)
	if err != nil {
		return err
	}
	// 部署 crd 和 controller
}
