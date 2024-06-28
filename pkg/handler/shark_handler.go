package handler

import (
	"github.com/985492783/sparrow-go/integration"
	server2 "github.com/985492783/sparrow-go/pkg/server"
	"github.com/985492783/sparrow-go/pkg/utils"
)

type SharkHandler struct {
}

type SharkRequest struct {
	integration.Metadata
}
type SharkResponse struct {
}

func (s *SharkResponse) Code() int {
	return 200
}

func (s *SharkResponse) SetCode(i int) {
}

func (shark *SharkHandler) Handler(request integration.Request) integration.Response {
	return &SharkResponse{}
}

func (shark *SharkHandler) GetType() (string, utils.Construct) {
	return utils.GetType(SharkRequest{}), func() integration.Request {
		return &SharkRequest{}
	}
}

func (shark *SharkHandler) GetPermit(request integration.Request) string {
	return ""
}

func NewSharkHandler() *SharkHandler {
	return &SharkHandler{}
}

var _ server2.RequestHandler = (*SharkHandler)(nil)
