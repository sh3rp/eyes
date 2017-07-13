package controller

import (
	"log"
	"net"

	"github.com/golang/protobuf/proto"
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
	pa.Connection.Write(data)
	return nil
}

type ProbeController struct {
	Agents map[string]*ProbeAgent
}

func NewProbeController() *ProbeController {
	return &ProbeController{
		Agents: make(map[string]*ProbeAgent),
	}
}

func (c *ProbeController) Start() {
	ln, err := net.Listen("tcp", ":12121")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go c.handle(conn)
	}
}

func (c *ProbeController) SendProbe(agentId string, cmd *messages.ProbeCommand) {
	c.Agents[agentId].SendCommand(cmd)
}

func (c *ProbeController) handle(conn net.Conn) {
	data := make([]byte, 4096)
	len, err := conn.Read(data)

	if err != nil {
		log.Printf("ERROR (read): %v", err)
		return
	}

	ack := &messages.ProbeACK{}
	err = proto.Unmarshal(data[:len], ack)

	if err != nil {
		log.Printf("ERROR (marshal): %v", err)
		return
	}

	c.Agents[ack.Id] = &ProbeAgent{
		Id:         ack.Id,
		Label:      ack.Label,
		Connection: conn,
	}

	log.Printf("Agent connected: %s (%v)", ack.Id, c.Agents[ack.Id])
}
