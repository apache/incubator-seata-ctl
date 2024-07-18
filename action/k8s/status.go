package k8s

import "github.com/spf13/cobra"

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "show seata status in kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		println("create config success!")
	},
}
