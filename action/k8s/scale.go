package k8s

import "github.com/spf13/cobra"

var ScaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "scale seata in kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		println("create config success!")
	},
}
