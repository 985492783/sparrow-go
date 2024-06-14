package center

import (
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/remote/server"
	"github.com/985492783/sparrow-go/pkg/util"
)

const (
	REGISTRY = "registry"
)

type SwitcherRequest struct {
	Name      string
	Kind      string
	AppName   string `validate:"required"`
	Ip        string `json:"ip"`
	ClassName string `json:"className"`
	SwitchMap map[string]SwitcherItem
}

type SwitcherItem struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
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

func (s *SwitcherHandler) GetType() (string, util.Construct) {
	return util.GetType(SwitcherRequest{}), func() integration.Request {
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
	return nil
}
