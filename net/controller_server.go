package net

import (
	"net"

	"github.com/sh3rp/eyes/controller"
	"github.com/sh3rp/eyes/db"
	"github.com/sh3rp/eyes/msg"
)

type ControllerServer struct {
	connection Connection
	controller controller.Controller
	db         db.EyesDB
}

func NewControllerServer(c net.Conn, ctrl controller.Controller, db db.EyesDB, info msg.NodeInfo) ControllerServer {
	conn := NewConnection(c, info, 10000)

	controller := ControllerServer{
		connection: conn,
		controller: ctrl,
		db:         db,
	}
	conn.SetHandler(controller)

	return controller
}

func (c ControllerServer) HandlePacket(pkt msg.Packet) {
}

func (c ControllerServer) HandleError(data []byte, err error) {
}
