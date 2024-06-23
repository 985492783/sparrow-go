package utils

import (
	"fmt"
	"github.com/985492783/sparrow-go/integration"
	"testing"
)

type SimpleResponse struct {
	Second     int
	StatusCode int
}

func (s *SimpleResponse) Code() int {
	return s.StatusCode
}

type SimpleRequest struct {
	Name string
}

func TestConvert(t *testing.T) {
	simple := &SimpleResponse{
		Second:     100,
		StatusCode: 200,
	}
	response, err := ConvertResponse(simple)
	if err != nil {
		return
	}
	fmt.Println(response)
}
func TestRequest(t *testing.T) {
	simple := &SimpleRequest{
		Name: "hahaha",
	}
	request, err := ConvertRequest(simple)
	if err != nil {
		return
	}
	fmt.Println(request)
	RegistryConstruct("utils.SimpleRequest", func() integration.Request {
		return &SimpleRequest{}
	})
	parseRequest, err := ParseRequest(request)
	if err != nil {
		return
	}
	fmt.Println(parseRequest)
}
