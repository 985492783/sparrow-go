package command

import (
	"context"
	"fmt"
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/handler"
	"github.com/985492783/sparrow-go/pkg/utils"
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
		queryCmd.AddCommand(getClassCmd())
	}
	switcherCmd.AddCommand(queryCmd)
}

func getNsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ns",
		Short: "Get Namespace",
		Run: func(cmd *cobra.Command, args []string) {
			if !isConnected {
				fmt.Println("please client first: use connect")
				return
			}
			request := &handler.SwitcherRequest{
				Kind: handler.QUERY,
				SwitcherQuery: &handler.SwitcherQuery{
					Level: "ns",
				},
			}
			response, err := queryNs(request)
			if err != nil {
				fmt.Printf("query err: %v\n", err)
				return
			}
			fmt.Printf("%v\n", *response)
		},
	}
}

func getClassCmd() *cobra.Command {
	class := &cobra.Command{
		Use:   "class",
		Short: "Get ClassMap",
		Run: func(cmd *cobra.Command, args []string) {
			if !isConnected {
				fmt.Println("please client first: use connect")
				return
			}
			ns, err := cmd.Flags().GetString("namespace")
			if err != nil {
				fmt.Printf("get namespace err: %v\n", err)
				return
			}
			request := &handler.SwitcherRequest{
				Kind: handler.QUERY,
				SwitcherQuery: &handler.SwitcherQuery{
					Level: "class",
				},
				Metadata: &integration.Metadata{
					NameSpace: ns,
				},
			}
			response, err := queryNs(request)
			if err != nil {
				fmt.Printf("query err: %v\n", err)
				return
			}
			fmt.Printf("%v\n", *response)
		},
	}
	class.Flags().String("namespace", "public", "query class map from [namespace]")
	return class
}

func query(request integration.Request) (integration.Response, error) {
	convertRequest, err := utils.ConvertRequest(request)
	if err != nil {
		return nil, err
	}
	payload, err := connect.req.Request(context.Background(), convertRequest)
	if err != nil {
		return nil, err
	}
	response, err := utils.ParseResponseByType(payload, &handler.SwitcherResponse{})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func queryNs(request integration.Request) (*string, error) {
	response, err := query(request)
	if err != nil {
		return nil, err
	}
	return &response.(*handler.SwitcherResponse).Resp, nil
}
