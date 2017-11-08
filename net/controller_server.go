package net

import (
	"net"

	"github.com/rs/zerolog/log"
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

func (c ControllerServer) HandleHello(m msg.Hello) {
	log.Error().Msg("error")
}
func (c ControllerServer) HandleKeepalive(m msg.KeepAlive) {
	log.Error().Msg("error")
}

func (c ControllerServer) HandleScheduleActionConfig(m msg.ScheduleActionConfig) {

}

func (c ControllerServer) HandleResult(m msg.Result) {
	log.Error().Msgf("Received result, should not have: %v", m)
}

func (c ControllerServer) HandleUnScheduleActionConfig(m msg.UnscheduleActionConfig) {
}

func (c ControllerServer) HandleRunActionConfig(m msg.RunActionConfig) {
	log.Error().Msg("error")
}

func (c ControllerServer) HandleAllActionConfigs(m msg.AllActionConfigs) {
}
