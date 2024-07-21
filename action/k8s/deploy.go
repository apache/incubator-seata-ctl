package k8s

import (
	"github.com/spf13/cobra"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy seata in kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		deploy()
	},
}

var version string

func init() {
	DeployCmd.PersistentFlags().StringVar(&version, "version", "lastest", "crd version")
}

func deploy() error {
	//获取动态kubeclient
	//client, err := common.GetClient()
	//if err != nil {
	//	return err
	//}
	//操作deploy
	//deployment := &appV1.Deployment{
	//	ObjectMeta: metaV1.ObjectMeta{
	//		Name: "nginx",
	//		Labels: map[string]string{
	//			"app": "nginx",
	//		},
	//		Namespace: "default",
	//	},
	//	Spec: appV1.DeploymentSpec{},
	//}
	//deployment, err = client.AppsV1().Deployments("default").Create(context.Background(), deployment, metaV1.CreateOptions{})
	return nil
}
