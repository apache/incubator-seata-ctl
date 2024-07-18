package k8s

import (
	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install k8s crd controller",
	Run: func(cmd *cobra.Command, args []string) {
		println("create config success!")
	},
}
