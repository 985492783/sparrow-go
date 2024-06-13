package util

import (
	"encoding/json"
	"fmt"
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/remote/pb"
	"google.golang.org/protobuf/types/known/anypb"
	"reflect"
	"strconv"
	"sync"
)

type Construct func() integration.Request

var structMap sync.Map

func RegistryConstruct(t string, construct Construct) {
	structMap.Store(t, construct)
}
func ConvertResponse(response integration.Response) (*pb.Payload, error) {
	marshal, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	realType := GetType(response)
	meta := &pb.Metadata{
		Type:    realType,
		Headers: make(map[string]string, 4),
	}
	meta.GetHeaders()[integration.StatusCode] = strconv.Itoa(response.Code())
	payload := &pb.Payload{
		Metadata: meta,
		Body: &anypb.Any{
			TypeUrl: "type.googleapis.com/json",
			Value:   marshal,
		},
	}
	return payload, nil
}

func ConvertRequest(request integration.Request) (*pb.Payload, error) {
	marshal, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	realType := GetType(request)
	meta := &pb.Metadata{Type: realType}
	payload := &pb.Payload{
		Metadata: meta,
		Body: &anypb.Any{
			TypeUrl: "type.googleapis.com/json",
			Value:   marshal,
		},
	}
	return payload, nil
}

func ParseRequest(payload *pb.Payload) (integration.Request, error) {
	t := payload.Metadata.Type
	if fun, ok := structMap.Load(t); ok {
		obj := fun.(Construct)()
		err := json.Unmarshal(payload.Body.Value, obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
	return nil, fmt.Errorf("not found struct type: %s", t)
}

func GetType(request interface{}) string {
	tt := reflect.TypeOf(request)
	if tt.Kind() == reflect.Ptr {
		return tt.Elem().String()
	}
	return tt.String()
}
