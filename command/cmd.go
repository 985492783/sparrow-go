package command

import (
	"errors"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"io"
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
	rootCmd.AddCommand(&cobra.Command{
		Use:   "exit",
		Short: "exit command",
		Run: func(cmd *cobra.Command, args []string) {
			disConnect()
			os.Exit(0)
		},
	})
}
func Execute() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "sparrow> ",
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Println("Error initializing readline:", err)
		return
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // 处理错误或退出
			if errors.Is(err, readline.ErrInterrupt) {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}
			fmt.Println("Error reading line:", err)
			continue
		}

		// 去除输入的换行符和前后空格
		line = strings.TrimSpace(line)

		rootCmd.SetArgs(strings.Split(line, " "))
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
		}
	}

}
