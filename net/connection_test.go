package net

import (
	"net"
	"testing"

	"github.com/sh3rp/eyes/msg"
	"github.com/stretchr/testify/assert"
)

func TestSimplePacket(t *testing.T) {
	c1, c2 := net.Pipe()
	cHandler := &mockHandler{}
	aHandler := &mockHandler{}
	controller := NewConnection(c1, cHandler, 0)
	agent := NewConnection(c2, aHandler, 0)

	err := agent.Send(msg.Packet{Sender: msg.Packet_AGENT})
	assert.Nil(t, err)
	err = controller.Send(msg.Packet{Sender: msg.Packet_CONTROLLER})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(aHandler.packets))
	assert.Equal(t, 1, len(cHandler.packets))
}

func TestMalformedPacket(t *testing.T) {
	c1, c2 := net.Pipe()
	mockHandler := &mockHandler{}
	NewConnection(c1, mockHandler, 0)

	c2.Write([]byte("bogus"))

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
