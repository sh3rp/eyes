package controller

import (
	"bytes"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/messages"
)

type ProbeAgent struct {
	Id         string
	Label      string
	Connection net.Conn
}

func (pa *ProbeAgent) SendCommand(cmd *messages.ProbeCommand) error {
	data, err := proto.Marshal(cmd)

	if err != nil {
		return err
	}
	log.Info().Msgf("Writing command: %v", cmd)
	pa.Connection.Write(data)
	return nil
}

func (pa *ProbeAgent) ReadLoop(resultChannel chan *messages.ProbeResult) {
	for {
		data := make([]byte, 4096)
		len, err := pa.Connection.Read(data)

		if err != nil {
			log.Error().Msgf("ERROR (readLoop): %v", err)
		}

		ack := &messages.ProbeACK{}
		err = proto.Unmarshal(data[:len], ack)

		if err != nil {
			log.Error().Msgf("ERROR (unmarshal): %v", err)
		}

		switch ack.Type {
		case messages.ProbeACK_RESULT:
			resultChannel <- ack.Result
		}
	}
}

type ProbeController struct {
	Agents        map[string]*ProbeAgent
	ResultChannel chan *messages.ProbeResult
}

func NewProbeController() *ProbeController {
	return &ProbeController{
		Agents:        make(map[string]*ProbeAgent),
		ResultChannel: make(chan *messages.ProbeResult, 10),
	}
}

func (c *ProbeController) ResultReadLoop() {
	for {
		result := <-c.ResultChannel
		switch result.Type {
		case messages.ProbeResult_NOOP:
			cmp := bytes.Compare([]byte{0, 1, 2, 3, 4, 5, 6, 7}, result.Data)
			log.Info().Msgf("Agent %s probe test: %v", result.ProbeId, cmp == 0)
		}
	}
}

func (c *ProbeController) Start() {
	ln, err := net.Listen("tcp", ":12121")

	if err != nil {
		log.Error().Msgf("%v", err)
	}

	go c.ResultReadLoop()

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Error().Msgf("%v", err)
		}

		go c.handle(conn)
	}
}

func (c *ProbeController) SendProbe(agentId string, cmd *messages.ProbeCommand) {
	log.Info().Msgf("SendProbe: %s", cmd.Type)
	if v, ok := c.Agents[agentId]; ok {
		v.SendCommand(cmd)
	} else {
		log.Error().Msgf("SendProbe failed, no such agentId %s", agentId)
	}
}

func (c *ProbeController) TestProbe(agentId string) {
	c.SendProbe(agentId, &messages.ProbeCommand{
		Type: messages.ProbeCommand_NOOP,
		Host: "127.0.0.1",
	})
}

func (c *ProbeController) handle(conn net.Conn) {
	data := make([]byte, 4096)
	len, err := conn.Read(data)

	if err != nil {
		log.Error().Msgf("ERROR (read): %v", err)
		return
	}

	ack := &messages.ProbeACK{}
	err = proto.Unmarshal(data[:len], ack)

	if err != nil {
		log.Error().Msgf("ERROR (marshal): %v", err)
		return
	}

	c.Agents[ack.Id] = &ProbeAgent{
		Id:         ack.Id,
		Label:      ack.Label,
		Connection: conn,
	}
	go c.Agents[ack.Id].ReadLoop(c.ResultChannel)

	log.Info().Msgf("Agent connected: %s (%v)", ack.Id, c.Agents[ack.Id])
}
