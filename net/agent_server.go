package net

import (
	"net"

	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/msg"
	"github.com/sh3rp/eyes/util"
)

type AgentServer struct {
	connection Connection
	agent      agent.Agent
}

func NewAgentServer(c net.Conn, agent agent.Agent) AgentServer {
	connection := NewConnection(c, msg.NodeInfo{}, 10000)
	// TODO: fix this circular dependency
	agentServer := AgentServer{
		connection: connection,
		agent:      agent,
	}
	connection.SetHandler(agentServer)
	agent.AddResultHandler("agent", agentServer.shipResult)
	return agentServer
}

func (s AgentServer) HandlePacket(pkt msg.Packet) {
}

func (s AgentServer) HandleError(data []byte, err error) {
	log.Error().Msgf("error: %v (DATA: %v)", err, data)
}

func (s AgentServer) HandleHello(m msg.Hello) {
	log.Error().Msg("error")
}
func (s AgentServer) HandleKeepalive(m msg.KeepAlive) {
	log.Error().Msg("error")
}

func (s AgentServer) HandleScheduleActionConfig(m msg.ScheduleActionConfig) {
	c := PBtoAgentConfig(m.Config)
	s.agent.StoreActionConfig(c)
	s.agent.ScheduleAction(c.Id, m.Schedule)
}

func (s AgentServer) HandleResult(m msg.Result) {
	log.Error().Msgf("Received result, should not have: %v", m)
}

func (s AgentServer) HandleUnScheduleActionConfig(m msg.UnscheduleActionConfig) {
	c := PBtoAgentConfig(m.Config)
	s.agent.UnscheduleAction(c.Id)
}

func (s AgentServer) HandleRunActionConfig(m msg.RunActionConfig) {
	log.Error().Msg("error")
}

func (s AgentServer) HandleAllActionConfigs(m msg.AllActionConfigs) {
}

func (s AgentServer) shipResult(r agent.Result) {
	pb := ResultToPB(r)
	s.connection.Send(msg.Packet{
		Sender: msg.Packet_AGENT,
		Packet: &msg.Packet_Result{pb},
		Code:   0,
		Msg:    "ok",
	})
}

func ResultToPB(r agent.Result) *msg.Result {
	return &msg.Result{
		Id:        string(r.Id),
		ConfigId:  string(r.ConfigId),
		DataCode:  int32(r.DataCode),
		Data:      r.Data,
		Timestamp: r.Timestamp,
		Tags:      r.Tags,
	}
}

func PBtoAgentConfig(cfg *msg.ActionConfig) agent.ActionConfig {
	return agent.ActionConfig{
		Id:         util.ID(cfg.Id),
		Action:     int(cfg.Action),
		Parameters: cfg.Parameters,
	}
}
