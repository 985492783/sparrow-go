package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

var switcherCmd = &cobra.Command{
	Use:   "switcher",
	Short: "switcher engine",
}

func init() {
	queryCmd := &cobra.Command{
		Use:   "get",
		Short: "get ns|class|fields",
	}
	{
		queryCmd.AddCommand(getNsCmd())
	}
	switcherCmd.AddCommand(queryCmd)
}

func getNsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ns",
		Short: "Get Namespace",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("123")
		},
	}
}
