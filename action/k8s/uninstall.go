package k8s

import "github.com/spf13/cobra"

var UnDeployCmd = &cobra.Command{
	Use:   "undeploy",
	Short: "undeploy seata in kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		println("create config success!")
	},
}
