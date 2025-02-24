package server

import (
	"github.com/985492783/sparrow-go/integration"
	"github.com/985492783/sparrow-go/pkg/core"
	"github.com/985492783/sparrow-go/pkg/remote/pb"
	"github.com/985492783/sparrow-go/pkg/utils"
	"io"
	"log"
	"sync"
)

type RegistryRequest struct {
	Ip string `json:"ip"`
	*integration.Metadata
}
type RegistryResponse struct {
}

func (response *RegistryResponse) Code() int {
	return 200
}

func (response *RegistryResponse) SetCode(i int) {
}

type RequestServerStream struct {
	pb.UnsafeBiRequestStreamServer
	streamMap map[string]pb.BiRequestStream_RequestBiStreamServer
	lock      sync.RWMutex
	manager   *core.SwitcherManager
}

var stream *RequestServerStream

func NewRequestStream(manager *core.SwitcherManager) *RequestServerStream {
	if stream == nil {
		utils.RegistryConstruct(utils.GetType(RegistryRequest{}), func() integration.Request {
			return &RegistryRequest{
				Metadata: &integration.Metadata{},
			}
		})
		stream = &RequestServerStream{
			streamMap: make(map[string]pb.BiRequestStream_RequestBiStreamServer),
			manager:   manager,
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
				server.addStream(clientId, stream)
				open = true
				response, err := utils.ConvertResponse(&RegistryResponse{})
				if err != nil {
					continue
				}
				stream.Send(response)
			}
		}
		log.Println("deregister")
		server.manager.DeRegister(clientId)
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
