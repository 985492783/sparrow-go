package center

import (
	"fmt"
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/remote/server"
	"github.com/985492783/sparrow-go/pkg/util"
)

type SwitcherRequest struct {
	Name string
}

type SwitcherResponse struct {
	Resp       string
	statusCode int
}

func (c *SwitcherResponse) Code() int {
	return c.statusCode
}

type SwitcherHandler struct {
}

func NewSwitcherHandler() server.RequestHandler {
	return &SwitcherHandler{}
}

func (s *SwitcherHandler) Handler(payload integration.Request) integration.Response {
	request := payload.(*SwitcherRequest)
	fmt.Printf("switcher handler, %v\n", request)
	return &SwitcherResponse{
		Resp:       "success " + request.Name,
		statusCode: 200,
	}
}

func (s *SwitcherHandler) GetType() (string, util.Construct) {
	return util.GetType(SwitcherRequest{}), func() integration.Request {
		return &SwitcherRequest{}
	}
}

var _ server.RequestHandler = (*SwitcherHandler)(nil)
