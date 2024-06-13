package server

import (
	"context"
	"fmt"
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/remote/pb"
	"github.com/985492783/sparrow-go/pkg/util"
)

type RequestService struct {
	pb.UnimplementedRequestServer
}

var _ pb.RequestServer = (*RequestService)(nil)

func NewRequestService() *RequestService {
	return &RequestService{}
}
func (service *RequestService) Request(context context.Context, payload *pb.Payload) (*pb.Payload, error) {
	request, err := util.ParseRequest(payload)
	if err != nil {
		return nil, err
	}
	fmt.Println(request)
	resp := &integration.SuccessResponse{}
	response, err := util.ConvertResponse(resp)
	if err != nil {
		return nil, err
	}
	return response, nil
}
