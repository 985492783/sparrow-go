package server

import (
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/remote/pb"
	"github.com/985492783/sparrow-go/pkg/utils"
	"io"
	"sync"
)

type RegistryRequest struct {
	ClientId string `validate:"clientId"`
	Ip       string `json:"ip"`
}

type RequestServerStream struct {
	pb.UnsafeBiRequestStreamServer
	streamMap map[string]pb.BiRequestStream_RequestBiStreamServer
	lock      sync.RWMutex
}

var stream *RequestServerStream

func NewRequestStream() *RequestServerStream {
	if stream == nil {
		utils.RegistryConstruct(utils.GetType(RegistryRequest{}), func() integration.Request {
			return &RegistryRequest{}
		})
		stream = &RequestServerStream{
			streamMap: make(map[string]pb.BiRequestStream_RequestBiStreamServer),
		}
	}
	return stream
}
func (server *RequestServerStream) RequestBiStream(stream pb.BiRequestStream_RequestBiStreamServer) error {
	ch := make(chan interface{})
	go func() {
		var clientId string
		var open = false
		defer func() {
			if open {
				server.removeStream(clientId)
			}
			ch <- nil
			close(ch)
		}()

		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
			if !open && req.Metadata.Type == utils.GetType(RegistryRequest{}) {
				request, err := utils.ParseRequest(req)
				if err != nil {
					continue
				}
				registryRequest := request.(*RegistryRequest)
				clientId = registryRequest.ClientId
				server.addStream(registryRequest.ClientId, stream)
				open = true
			}

		}
	}()

	<-ch
	return nil
}

func (server *RequestServerStream) GetStreams() map[string]pb.BiRequestStream_RequestBiStreamServer {
	return server.streamMap
}

func (server *RequestServerStream) removeStream(clientId string) {
	server.lock.Lock()
	defer server.lock.Unlock()
	delete(server.streamMap, clientId)
}

func (server *RequestServerStream) addStream(clientId string, stream pb.BiRequestStream_RequestBiStreamServer) {
	server.lock.Lock()
	defer server.lock.Unlock()
	server.streamMap[clientId] = stream
}