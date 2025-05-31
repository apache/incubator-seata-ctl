package login

import (
	"fmt"
	"os"

	"github.com/seata/seata-ctl/seata"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Address string

// LoginCmd 定义 login 子命令
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Seata server",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Attempting to login...")
		Address = seata.GetAuth().GetAddress()
		err := seata.GetAuth().Login()
		if err != nil {
			fmt.Println("Login failed!")
			os.Exit(1)
		}
		fmt.Printf("Login successful to address: %s\n", Address)
		printPrompt(Address)
	},
}

func init() {
	credential := seata.GetAuth()
	LoginCmd.PersistentFlags().StringVar(&credential.ServerIP, "ip", "127.0.0.1", "Seata Server IP")
	LoginCmd.PersistentFlags().IntVar(&credential.ServerPort, "port", 7091, "Seata Server Admin Port")
	LoginCmd.PersistentFlags().StringVar(&credential.Username, "username", "seata", "Username")
	LoginCmd.PersistentFlags().StringVar(&credential.Password, "password", "seata", "Password")
	viper.BindPFlag("ip", LoginCmd.PersistentFlags().Lookup("ip"))
	viper.BindPFlag("port", LoginCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("username", LoginCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", LoginCmd.PersistentFlags().Lookup("password"))
}

func printPrompt(address string) {
	fmt.Printf("%s > ", address)
}
