package controller

import (
	"net"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/messages"
	"github.com/sh3rp/eyes/util"
)

type ProbeController struct {
	Agents            map[string]*ProbeAgent
	ResultChannel     chan *messages.AgentProbeResult
	ResultListeners   []func(*messages.AgentProbeResult)
	DisconnectChannel chan string
	agentLock         *sync.Mutex
}

func NewProbeController() *ProbeController {
	return &ProbeController{
		Agents:            make(map[string]*ProbeAgent),
		ResultChannel:     make(chan *messages.AgentProbeResult, 10),
		DisconnectChannel: make(chan string, 5),
		agentLock:         new(sync.Mutex),
	}
}

func (c *ProbeController) AddResultListener(f func(*messages.AgentProbeResult)) {
	c.ResultListeners = append(c.ResultListeners, f)
}

func (c *ProbeController) ResultReadLoop() {
	log.Info().Msgf("ResultReadLoop: starting")
	for {
		result := <-c.ResultChannel
		for _, listener := range c.ResultListeners {
			listener(result)
		}
	}
}

func (c *ProbeController) DisconnectHandler() {
	log.Info().Msgf("DisconnectHandler: starting")
	for {
		disconnect := <-c.DisconnectChannel
		c.agentLock.Lock()
		if _, ok := c.Agents[disconnect]; ok {
			delete(c.Agents, disconnect)
		}
		c.agentLock.Unlock()
	}
}

func (c *ProbeController) Start() {
	log.Info().Msgf("Controller: starting")

	ln, err := net.Listen("tcp", ":12121")

	if err != nil {
		log.Error().Msgf("Error listening on socket: %v", err)
		return
	}

	go c.ResultReadLoop()
	go c.DisconnectHandler()

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Error().Msgf("Error accepting connection: %v", err)
		} else {
			go c.handle(conn)
		}
	}
}

func (c *ProbeController) SendProbe(agentId string, latencyRequest *messages.ControllerLatencyRequest) string {
	if latencyRequest.ResultId == "" {
		latencyRequest.ResultId = util.GenID()
	}
	c.agentLock.Lock()
	defer c.agentLock.Unlock()
	if v, ok := c.Agents[agentId]; ok {
		v.SendCommand(&messages.ControllerMessage{
			Type:           messages.ControllerMessage_LATENCY_REQUEST,
			LatencyRequest: latencyRequest,
		})
	} else {
		log.Error().Msgf("SendProbe failed, no such agentId %s", agentId)
	}
	return latencyRequest.ResultId
}

func (c *ProbeController) handle(conn net.Conn) {
	data := make([]byte, 4096)
	len, err := conn.Read(data)

	if err != nil {
		log.Error().Msgf("ERROR handle (read): %v", err)
		return
	}

	agentMessage := &messages.AgentMessage{}
	err = proto.Unmarshal(data[:len], agentMessage)

	if err != nil {
		log.Error().Msgf("ERROR handle (marshal): %v", err)
		return
	}

	c.agentLock.Lock()

	c.Agents[agentMessage.Id] = &ProbeAgent{
		Id:         agentMessage.Id,
		Info:       agentMessage.Info,
		Connection: conn,
	}

	c.agentLock.Unlock()

	go c.Agents[agentMessage.Id].ReadLoop(c.ResultChannel, c.DisconnectChannel)

	log.Info().Msgf("Agent connected: %s (%v) - (%v)", agentMessage.Id, c.Agents[agentMessage.Id].Info.Ipaddress, c.Agents[agentMessage.Id].Info.Label)
}
