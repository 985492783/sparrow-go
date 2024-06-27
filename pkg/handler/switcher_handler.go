package handler

import (
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/core"
	"github.com/985492783/sparrow-go/pkg/db"
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
	ClassMap   map[string]map[string]*core.SwitcherItem
}

func (c *SwitcherResponse) Code() int {
	return c.statusCode
}

type SwitcherHandler struct {
	stream          *server2.RequestServerStream
	switcherManager *core.SwitcherManager
	database        *db.Database
}

func (s *SwitcherHandler) GetPermit(payload integration.Request) string {
	request := payload.(*SwitcherRequest)
	switch request.Kind {
	case REGISTRY:
		return utils.AuthSwitcherRegister
	}
	return ""
}

var _ server2.RequestHandler = (*SwitcherHandler)(nil)

func (s *SwitcherHandler) GetType() (string, utils.Construct) {
	return utils.GetType(SwitcherRequest{}), func() integration.Request {
		return &SwitcherRequest{}
	}
}

func NewSwitcherHandler(database *db.Database, switcherManager *core.SwitcherManager, stream *server2.RequestServerStream) server2.RequestHandler {
	return &SwitcherHandler{
		stream:          stream,
		database:        database,
		switcherManager: switcherManager,
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
	// TODO 从底层获取持久化配置
	classMap := request.ClassMap
	mp := make(map[string]map[string]*core.SwitcherItem)

	for class, fieldMap := range classMap {
		properties := s.database.GetData(request.NameSpace, "switcher", request.AppName+"@@"+class)
		for field, item := range fieldMap {
			if fieldData, ok := properties.Get(field); ok && utils.IsTypeOf(fieldData, item.Type) {
				item.Value = fieldData
				clm, ok := mp[class]
				if !ok {
					clm = make(map[string]*core.SwitcherItem)
					mp[class] = clm
				}
				clm[field] = item
			}
		}
	}
	go s.switcherManager.Register(request.ClientId, request.NameSpace, request.AppName, request.Ip, classMap)
	return &SwitcherResponse{
		statusCode: 200,
		ClassMap:   mp,
	}
}
