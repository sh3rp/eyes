package net

import (
	"net"
	"testing"
	"time"

	"github.com/sh3rp/eyes/msg"
	"github.com/sh3rp/eyes/util"
	"github.com/stretchr/testify/assert"
)

var INFO = util.GenNodeInfo(util.NewId())

func TestSimplePacket(t *testing.T) {
	c1, c2 := net.Pipe()
	cHandler := &mockHandler{}
	aHandler := &mockHandler{}
	controller := NewConnection(c1, INFO, 600)
	controller.SetHandler(cHandler)
	agent := NewConnection(c2, INFO, 600)
	agent.SetHandler(aHandler)
	agent.Send(msg.Packet{Sender: msg.Packet_AGENT, Packet: &msg.Packet_Hello{&msg.Hello{}}})
	controller.Send(msg.Packet{Sender: msg.Packet_CONTROLLER, Packet: &msg.Packet_Hello{&msg.Hello{}}})
	time.Sleep(1 * time.Second)
	assert.Equal(t, 2, len(aHandler.packets))
	assert.Equal(t, 2, len(cHandler.packets))
	assert.Equal(t, 1, cHandler.keepalives)
	assert.Equal(t, 1, aHandler.keepalives)
}

func TestMalformedPacket(t *testing.T) {
	c1, c2 := net.Pipe()
	mockHandler := &mockHandler{}
	c := NewConnection(c1, INFO, 10000)
	c.SetHandler(mockHandler)

	c2.Write([]byte("bogus"))

	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, len(mockHandler.errors))
}

type mockHandler struct {
	packets    []msg.Packet
	errors     []dataError
	keepalives int
}

type dataError struct {
	data []byte
	err  error
}

func (h *mockHandler) HandlePacket(pkt msg.Packet) {
	if pkt.GetKeepalive() != nil {
		h.keepalives = h.keepalives + 1
	}
	h.packets = append(h.packets, pkt)
}

func (h *mockHandler) HandleError(data []byte, err error) {
	h.errors = append(h.errors, dataError{data, err})
}

func (h *mockHandler) HandleHello(m msg.Hello) {
}
func (h *mockHandler) HandleKeepalive(m msg.KeepAlive) {
}
func (h *mockHandler) HandleScheduleActionConfig(m msg.ScheduleActionConfig) {
}
func (h *mockHandler) HandleResult(m msg.Result) {
}
func (h *mockHandler) HandleUnScheduleActionConfig(m msg.UnscheduleActionConfig) {
}
func (h *mockHandler) HandleRunActionConfig(m msg.RunActionConfig) {
}
func (h *mockHandler) HandleAllActionConfigs(msg.AllActionConfigs) {

}
