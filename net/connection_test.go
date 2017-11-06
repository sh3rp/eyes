package net

import (
	"net"
	"testing"
	"time"

	"github.com/sh3rp/eyes/msg"
	"github.com/stretchr/testify/assert"
)

var INFO = msg.NodeInfo{
	Id: "dummy",
}

func TestSimplePacket(t *testing.T) {
	c1, c2 := net.Pipe()
	cHandler := &mockHandler{}
	aHandler := &mockHandler{}
	controller := NewConnection(c1, INFO)
	controller.SetHandler(cHandler)
	agent := NewConnection(c2, INFO)
	agent.SetHandler(aHandler)
	agent.Send(msg.Packet{Sender: msg.Packet_AGENT})
	controller.Send(msg.Packet{Sender: msg.Packet_CONTROLLER})
	time.Sleep(1 * time.Second)
	assert.Equal(t, 1, len(aHandler.packets))
	assert.Equal(t, 1, len(cHandler.packets))
}

func TestMalformedPacket(t *testing.T) {
	c1, c2 := net.Pipe()
	mockHandler := &mockHandler{}
	c := NewConnection(c1, INFO)
	c.SetHandler(mockHandler)

	c2.Write([]byte("bogus"))

	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, len(mockHandler.errors))
}

type mockHandler struct {
	packets []msg.Packet
	errors  []dataError
}

type dataError struct {
	data []byte
	err  error
}

func (h *mockHandler) HandlePacket(pkt msg.Packet) {
	h.packets = append(h.packets, pkt)
}

func (h *mockHandler) HandleError(data []byte, err error) {
	h.errors = append(h.errors, dataError{data, err})
}
