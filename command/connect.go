package command

import (
	"github.com/spf13/cobra"
)

var (
	username    string
	password    string
	authEnabled = true
)
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect to grpc",
}

func init() {
	connectCmd.Flags().StringVarP(&username, "username", "u", "", "cli username")
	connectCmd.Flags().StringVarP(&password, "password", "p", "", "cli password")
	connectCmd.Flags().BoolVar(&authEnabled, "authEnabled", false, "cli auth disabled")
}
