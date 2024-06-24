package server

import (
	"context"
	"fmt"
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/remote/pb"
	"github.com/985492783/sparrow-go/pkg/utils"
)

type RequestHandler interface {
	Handler(request integration.Request) integration.Response

	GetType() (string, utils.Construct)
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
	utils.RegistryConstruct(name, construct)
	service.handlers[name] = handler
}
func (service *RequestService) Request(_ context.Context, payload *pb.Payload) (*pb.Payload, error) {
	//TODO 鉴权 clientId
	request, err := utils.ParseRequest(payload)
	if err != nil {
		return nil, err
	}
	realType := utils.GetType(request)
	handler, ok := service.handlers[realType]
	if !ok {
		return nil, fmt.Errorf("no handler for type %s", realType)
	}
	resp := handler.Handler(request)

	response, err := utils.ConvertResponse(resp)
	if err != nil {
		return nil, err
	}
	return response, nil
}
