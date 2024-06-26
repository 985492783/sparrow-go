package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/config"
	"github.com/985492783/sparrow-go/pkg/remote/pb"
	"github.com/985492783/sparrow-go/pkg/utils"
)

type RequestHandler interface {
	Handler(request integration.Request) integration.Response

	GetType() (string, utils.Construct)

	GetPermit(request integration.Request) string
}

type EmptyPermit struct {
}

func (p *EmptyPermit) GetPermit() string {
	return ""
}

const (
	username = "username"
	password = "password"
)

type RequestService struct {
	pb.UnimplementedRequestServer
	handlers map[string]RequestHandler
	config   *config.SparrowConfig
}

var _ pb.RequestServer = (*RequestService)(nil)

func NewRequestService(config *config.SparrowConfig) *RequestService {
	return &RequestService{
		handlers: make(map[string]RequestHandler),
		config:   config,
	}
}

func (service *RequestService) RegisterHandler(handler RequestHandler) {
	name, construct := handler.GetType()
	utils.RegistryConstruct(name, construct)
	service.handlers[name] = handler
}
func (service *RequestService) Request(_ context.Context, payload *pb.Payload) (*pb.Payload, error) {
	request, err := utils.ParseRequest(payload)
	if err != nil {
		return nil, err
	}
	realType := utils.GetType(request)
	handler, ok := service.handlers[realType]
	if !ok {
		return nil, fmt.Errorf("no handler for type %s", realType)
	}

	err = service.authPermit(request, handler.GetPermit(request))
	if err != nil {
		return nil, err
	}

	resp := handler.Handler(request)

	response, err := utils.ConvertResponse(resp)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (service *RequestService) authPermit(request integration.Request, permit string) error {
	if !service.config.AuthEnabled {
		return nil
	}
	user := request.GetHeader(username, "")
	pass := request.GetHeader(password, "")
	if user == "" {
		return errors.New("username is empty")
	}
	return service.config.Authority(user, pass, permit)
}
