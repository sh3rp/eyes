package net

import (
	"net"

	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/msg"
)

type AgentServer struct {
	connection Connection
	agent      agent.Agent
}

func NewAgentServer(c net.Conn, agent agent.Agent) AgentServer {
	connection := NewConnection(c)
	// TODO: fix this cirular dependency
	agentServer := AgentServer{
		connection: connection,
		agent:      agent,
	}
	connection.SetHandler(agentServer)
	agent.AddResultHandler("agent", agentServer.handleResult)
	return agentServer
}

func (s AgentServer) HandlePacket(pkt msg.Packet) {
	if pkt.Type == msg.Packet_PROBE {
		probe := pkt.GetProbe()
		switch probe.Action {
		case msg.Probe_ACTIVATE:
			c := PBtoAgentConfig(probe.Config)
			s.agent.StoreActionConfig(c)
			s.agent.ScheduleAction(c.Id, probe.Schedule)
		case msg.Probe_DEACTIVATE:
			c := PBtoAgentConfig(probe.Config)
			s.agent.UnscheduleAction(c.Id)
		case msg.Probe_RUN_ONCE:
		case msg.Probe_DELETE:
			c := PBtoAgentConfig(probe.Config)
			s.agent.UnscheduleAction(c.Id)
			s.agent.DeleteActionConfig(c.Id)
		}
	}
}

func (s AgentServer) HandleError(data []byte, err error) {
	log.Error().Msgf("error: %v (DATA: %v)", err, data)
}

func (s AgentServer) handleResult(r agent.Result) {
	pb := ResultToPB(r)
	s.connection.Send(msg.Packet{
		Sender:       msg.Packet_AGENT,
		Type:         msg.Packet_PROBE_RESULT,
		Packet:       &msg.Packet_Result{pb},
		ErrorCode:    0,
		ErrorMessage: "ok",
	})
}

func ResultToPB(r agent.Result) *msg.Result {
	return &msg.Result{
		Id:        r.Id,
		ConfigId:  r.ConfigId,
		DataCode:  int32(r.DataCode),
		Data:      r.Data,
		Timestamp: r.Timestamp,
		Tags:      r.Tags,
	}
}

func PBtoAgentConfig(cfg *msg.ActionConfig) agent.ActionConfig {
	return agent.ActionConfig{
		Id:         cfg.Id,
		Action:     int(cfg.Action),
		Parameters: cfg.Parameters,
	}
}
