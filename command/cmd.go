package command

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Sparrow CLI tool",
	Long:  "Sparrow is a light switcherCenter and logCenter, cli can easy update sparrow",
}

func init() {
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(switcherCmd)
}
func Execute() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		name, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		name = strings.TrimSpace(name)
		rootCmd.SetArgs(strings.Split(name, " "))
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
		}
	}

}
