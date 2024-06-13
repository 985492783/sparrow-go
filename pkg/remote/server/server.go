package server

import (
	"context"
	"fmt"
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/remote/pb"
	"github.com/985492783/sparrow-go/pkg/util"
)

type RequestHandler interface {
	Handler(request integration.Request) integration.Response

	GetType() (string, util.Construct)
}

type RequestService struct {
	pb.UnimplementedRequestServer
	handlers map[string]RequestHandler
}

var _ pb.RequestServer = (*RequestService)(nil)

func NewRequestService() *RequestService {
	return &RequestService{
		handlers: make(map[string]RequestHandler),
	}
}

func (service *RequestService) RegisterHandler(handler RequestHandler) {
	name, construct := handler.GetType()
	util.RegistryConstruct(name, construct)
	service.handlers[name] = handler
}
func (service *RequestService) Request(_ context.Context, payload *pb.Payload) (*pb.Payload, error) {
	request, err := util.ParseRequest(payload)
	if err != nil {
		return nil, err
	}
	realType := util.GetType(request)
	handler, ok := service.handlers[realType]
	if !ok {
		return nil, fmt.Errorf("no handler for type %s", realType)
	}
	resp := handler.Handler(request)

	response, err := util.ConvertResponse(resp)
	if err != nil {
		return nil, err
	}
	return response, nil
}
