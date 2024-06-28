package command

import (
	"context"
	"fmt"
	"github.com/985492783/sparrow-go/pkg/handler"
	"github.com/985492783/sparrow-go/pkg/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/connectivity"
	"strconv"
)

var (
	addr        = ":9854"
	username    string
	password    string
	isConnected = false
	connect     *client
)
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connectc to grpc",
	Run: func(cmd *cobra.Command, args []string) {
		if connect == nil || connect.conn.GetState() != connectivity.Ready {
			if connect != nil && connect.conn.GetState() != connectivity.Ready {
				fmt.Printf("connect is not ready\n")
			}
			con, err := newConnect(addr)
			if err != nil {
				fmt.Printf("new connect err: %v\n", err)
				isConnected = false
				return
			}
			connect = con
		}
		fmt.Printf("connect is ready\n")
		auth()
	},
}

func init() {
	connectCmd.Flags().StringVarP(&username, "username", "u", "", "cli username")
	connectCmd.Flags().StringVarP(&password, "password", "p", "", "cli password")
}

func auth() {
	request, err := utils.ConvertRequest(&handler.SharkRequest{})
	if err != nil {
		fmt.Printf("convert request err: %v\n", err)
		return
	}
	payload, err := connect.req.Request(context.Background(), request)
	if err != nil {
		fmt.Printf("connect request err: %v\n", err)
		return
	}
	code, ok := payload.Metadata.Headers["StatusCode"]
	if !ok {
		fmt.Printf("no StatusCode\n")
		isConnected = false
	} else if co, err := strconv.Atoi(code); err != nil || co != 200 {
		fmt.Printf("StatusCode: %d\n", co)
		isConnected = false
	} else {
		fmt.Printf("grpc sharkhand success\n")
		isConnected = true
	}
}

func disConnect() {
	if connect != nil {
		isConnected = false
		connect.conn.Close()
		connect = nil
	}
}
