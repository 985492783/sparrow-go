package handler

import (
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/core"
	server2 "github.com/985492783/sparrow-go/pkg/server"
	"github.com/985492783/sparrow-go/pkg/utils"
)

const (
	REGISTRY = "registry"
)

type SwitcherRequest struct {
	Kind     string
	AppName  string                                   `validate:"required"`
	Ip       string                                   `json:"ip"`
	ClassMap map[string]map[string]*core.SwitcherItem `json:"classMap"`
	integration.Metadata
}

type SwitcherResponse struct {
	Resp       string
	statusCode int
}

func (c *SwitcherResponse) Code() int {
	return c.statusCode
}

type SwitcherHandler struct {
	stream *server2.RequestServerStream
}

var _ server2.RequestHandler = (*SwitcherHandler)(nil)

func (s *SwitcherHandler) GetType() (string, utils.Construct) {
	return utils.GetType(SwitcherRequest{}), func() integration.Request {
		return &SwitcherRequest{}
	}
}

func NewSwitcherHandler(stream *server2.RequestServerStream) server2.RequestHandler {
	return &SwitcherHandler{
		stream: stream,
	}
}

func (s *SwitcherHandler) Handler(payload integration.Request) integration.Response {
	request := payload.(*SwitcherRequest)
	switch request.Kind {
	case REGISTRY:
		return s.registerHandler(request)
	}
	return &SwitcherResponse{
		Resp:       "Not Found " + request.Kind,
		statusCode: 404,
	}
}

func (s *SwitcherHandler) registerHandler(request *SwitcherRequest) *SwitcherResponse {
	if _, ok := s.stream.GetStreams()[request.ClientId]; !ok {
		return &SwitcherResponse{statusCode: 301, Resp: "clientId Not exist"}
	}
	core.Register(request.ClientId, request.NameSpace, request.AppName, request.Ip, request.ClassMap)
	return &SwitcherResponse{statusCode: 200}
}
