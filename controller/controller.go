package controller

import (
	"bytes"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/oklog/ulid"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/messages"
)

type ProbeController struct {
	Agents            map[string]*ProbeAgent
	ResultChannel     chan *messages.ProbeResult
	ResultListeners   []func(*messages.ProbeResult)
	DisconnectChannel chan string
	agentLock         *sync.Mutex
}

func NewProbeController() *ProbeController {
	return &ProbeController{
		Agents:            make(map[string]*ProbeAgent),
		ResultChannel:     make(chan *messages.ProbeResult, 10),
		DisconnectChannel: make(chan string, 5),
		agentLock:         new(sync.Mutex),
	}
}

func (c *ProbeController) AddResultListener(f func(*messages.ProbeResult)) {
	c.ResultListeners = append(c.ResultListeners, f)
}

func (c *ProbeController) ResultReadLoop() {
	log.Info().Msgf("ResultReadLoop: starting")
	for {
		result := <-c.ResultChannel
		switch result.Type {
		case messages.ProbeResult_NOOP:
			cmp := bytes.Compare([]byte{0, 1, 2, 3, 4, 5, 6, 7}, result.Data)
			log.Info().Msgf("Agent %s (%s) probe test: %v", result.ProbeId, result.Host, cmp == 0)
		}
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

func (c *ProbeController) SendProbe(agentId string, cmd *messages.ProbeCommand) string {
	log.Info().Msgf("SendProbe: %s", cmd.Type)
	cmd.Id = genID()
	c.agentLock.Lock()
	defer c.agentLock.Unlock()
	if v, ok := c.Agents[agentId]; ok {
		v.SendCommand(cmd)
	} else {
		log.Error().Msgf("SendProbe failed, no such agentId %s", agentId)
	}
	return cmd.Id
}

func (c *ProbeController) TestProbe(agentId string) string {
	return c.SendProbe(agentId, &messages.ProbeCommand{
		Type: messages.ProbeCommand_NOOP,
		Host: "127.0.0.1",
	})
}

func (c *ProbeController) handle(conn net.Conn) {
	data := make([]byte, 4096)
	len, err := conn.Read(data)

	if err != nil {
		log.Error().Msgf("ERROR handle (read): %v", err)
		return
	}

	ack := &messages.ProbeACK{}
	err = proto.Unmarshal(data[:len], ack)

	if err != nil {
		log.Error().Msgf("ERROR handle (marshal): %v", err)
		return
	}

	c.agentLock.Lock()

	c.Agents[ack.Id] = &ProbeAgent{
		Id:         ack.Id,
		IPAddress:  ack.Ipaddress,
		Label:      ack.Label,
		Location:   ack.Location,
		Connection: conn,
	}

	c.agentLock.Unlock()

	go c.Agents[ack.Id].ReadLoop(c.ResultChannel, c.DisconnectChannel)

	log.Info().Msgf("Agent connected: %s (%v)", ack.Id, c.Agents[ack.Id])
}

func genID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
