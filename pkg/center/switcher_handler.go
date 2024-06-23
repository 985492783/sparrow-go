package center

import (
	"fmt"
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/remote/server"
	"github.com/985492783/sparrow-go/pkg/utils"
)

const (
	REGISTRY = "registry"
)

type SwitcherRequest struct {
	Kind     string
	AppName  string                              `validate:"required"`
	Ip       string                              `json:"ip"`
	ClassMap map[string]map[string]*SwitcherItem `json:"classMap"`
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

var _ server.RequestHandler = (*SwitcherHandler)(nil)

func (s *SwitcherHandler) GetType() (string, utils.Construct) {
	return utils.GetType(SwitcherRequest{}), func() integration.Request {
		return &SwitcherRequest{}
	}
}

func NewSwitcherHandler() server.RequestHandler {
	return &SwitcherHandler{}
}

func (s *SwitcherHandler) Handler(payload integration.Request) integration.Response {
	request := payload.(*SwitcherRequest)
	switch request.Kind {
	case REGISTRY:
		return registerHandler(request)
	}
	return &SwitcherResponse{
		Resp:       "Not Found " + request.Kind,
		statusCode: 404,
	}
}

func registerHandler(request *SwitcherRequest) *SwitcherResponse {
	fmt.Println("registry")
	return &SwitcherResponse{statusCode: 200}
}
