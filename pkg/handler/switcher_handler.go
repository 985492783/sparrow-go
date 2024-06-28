package handler

import (
	"encoding/json"
	"fmt"
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/core"
	"github.com/985492783/sparrow-go/pkg/db"
	server2 "github.com/985492783/sparrow-go/pkg/server"
	"github.com/985492783/sparrow-go/pkg/utils"
)

const (
	REGISTRY = "registry"
	QUERY    = "query"
)

type SwitcherRequest struct {
	Kind     string                                   `json:"kind"`
	AppName  string                                   `json:"appName"`
	Ip       string                                   `json:"ip"`
	ClassMap map[string]map[string]*core.SwitcherItem `json:"classMap"`
	*integration.Metadata
	*SwitcherQuery
}

type SwitcherQuery struct {
	Level string `json:"level"`
}

type SwitcherResponse struct {
	*integration.ResponseData
	ClassMap map[string]map[string]*core.SwitcherItem `json:"classMap"`
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
	case QUERY:
		return utils.AuthSwitcherList
	}
	return ""
}

var _ server2.RequestHandler = (*SwitcherHandler)(nil)

func (s *SwitcherHandler) GetType() (string, utils.Construct) {
	return utils.GetType(SwitcherRequest{}), func() integration.Request {
		return &SwitcherRequest{
			Metadata: &integration.Metadata{},
		}
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
	case QUERY:
		return s.queryHandler(request)
	}
	return &SwitcherResponse{
		ResponseData: &integration.ResponseData{
			Resp:       "Not Found " + request.Kind,
			StatusCode: 404,
		},
	}
}

func (s *SwitcherHandler) registerHandler(request *SwitcherRequest) *SwitcherResponse {
	if _, ok := s.stream.GetStreams()[request.ClientId]; !ok {
		return &SwitcherResponse{
			ResponseData: &integration.ResponseData{
				StatusCode: 301,
				Resp:       "clientId Not exist",
			},
		}
	}

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
		ResponseData: &integration.ResponseData{
			StatusCode: 200,
		},
		ClassMap: mp,
	}
}

func (s *SwitcherHandler) queryHandler(request *SwitcherRequest) integration.Response {
	level := request.SwitcherQuery.Level
	ns := request.NameSpace
	response := &SwitcherResponse{
		ResponseData: &integration.ResponseData{
			StatusCode: 200,
		},
	}
	switch level {
	case "ns":
		response.Resp = fmt.Sprintf("%v", s.switcherManager.GetNs())
	case "class":
		data, err := json.Marshal(s.switcherManager.GetJSON(ns))
		if err != nil {
			response.StatusCode = 301
			response.Resp = err.Error()
		} else {
			response.Resp = string(data)
		}
	}
	return response
}
