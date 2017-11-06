package net

import (
	"net"

	"github.com/sh3rp/eyes/controller"
	"github.com/sh3rp/eyes/msg"
)

type AgentController struct {
	connection Connection
	controller controller.Controller
}

func NewAgentController(c net.Conn, ctrl controller.Controller) AgentController {
	conn := NewConnection(c, msg.NodeInfo{})

	controller := AgentController{
		connection: conn,
		controller: ctrl,
	}
	conn.SetHandler(controller)

	return controller
}

func (c AgentController) HandlePacket(pkt msg.Packet) {
}

func (c AgentController) HandleError(data []byte, err error) {
}
