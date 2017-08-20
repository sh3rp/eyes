package controller

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/sh3rp/eyes/messages"
)

type GRPCServer struct {
	Controller *ProbeController
	port       int
}

func NewGRPCServer(port int, controller *ProbeController) *GRPCServer {
	return &GRPCServer{
		Controller: controller,
		port:       port,
	}
}

func (s *GRPCServer) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		panic(err)
	}
	gServer := grpc.NewServer()
	messages.RegisterControllerServer(gServer, s)
	gServer.Serve(listener)
}

func (s *GRPCServer) GetControllerInfo(ctx context.Context, in *messages.Empty) (*messages.ControllerInfoResponse, error) {
	maj, min, pat := s.Controller.GetVersion()
	versionStr := fmt.Sprintf("v%d.%d.%d", maj, min, pat)
	return &messages.ControllerInfoResponse{
		Version:         versionStr,
		ConnectedAgents: int32(len(s.Controller.Agents)),
	}, nil
}

func (s *GRPCServer) CreateResultQueue(ctx context.Context, in *messages.CreateResultQueueRequest) (*messages.CreateResultQueueResponse, error) {
	return nil, nil
}

func (s *GRPCServer) GetAgents(ctx context.Context, req *messages.AgentRequest) (*messages.AgentResponse, error) {
	if s.Controller.Agents == nil || len(s.Controller.Agents) == 0 {
		return nil, errors.New("No agents connected to controller")
	}

	agentsMap := make(map[string]*messages.AgentInfo)

	for _, v := range s.Controller.Agents {
		agentsMap[v.Id] = v.Info
	}

	response := &messages.AgentResponse{
		Agents: agentsMap,
	}

	return response, nil
}

func (s *GRPCServer) SendProbe(ctx context.Context, req *messages.ProbeRequest) (*messages.ProbeResponse, error) {
	response := &messages.ProbeResponse{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	s.Controller.SendProbeCallback(req.AgentId, req.LatencyRequest, func(result *messages.AgentProbeResult) {
		response.Result = result
		wg.Done()
	})
	wg.Wait()
	return response, nil
}

func (s *GRPCServer) ScheduleProbe(ctx context.Context, req *messages.ScheduleProbeRequest) (*messages.ScheduleProbeResponse, error) {
	return nil, nil
}
