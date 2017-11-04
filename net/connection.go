package net

import (
	"errors"
	"net"

	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/msg"
)

var DEFAULT_BUFFER_SIZE = 8128

type Connection interface {
	Send(msg.Packet) error
	SetHandler(PacketHandler) error
}

type PacketHandler interface {
	HandlePacket(msg.Packet)
	HandleError([]byte, error)
}

type connection struct {
	conn       net.Conn
	bufferSize int
	handler    PacketHandler
}

func NewConnection(conn net.Conn) Connection {
	c := &connection{
		conn:       conn,
		bufferSize: DEFAULT_BUFFER_SIZE,
	}
	go c.read()

	return c
}

func (c *connection) SetHandler(handler PacketHandler) error {
	c.handler = handler
	return nil
}

func (c *connection) Send(pkt msg.Packet) error {
	return c.write(pkt)
}

func (c *connection) write(pkt msg.Packet) error {
	if c.conn == nil {
		return errors.New("Connection not set, cannot write packet")
	}
	data, err := proto.Marshal(&pkt)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *connection) read() {
	if c.conn == nil || c.handler == nil {
		if c.handler != nil {
			c.handler.HandleError(nil, errors.New("No connection set, cannot read"))
		} else {
			log.Error().Msgf("Catastrophic error: no handler set for inbound data")
		}
		return
	}
	for {
		data := make([]byte, c.bufferSize)
		len, err := c.conn.Read(data)

		if err == nil {
			pkt := msg.Packet{}
			marshalError := proto.Unmarshal(data[:len], &pkt)

			if marshalError != nil {
				c.handler.HandleError(data, marshalError)
			} else {
				c.handler.HandlePacket(pkt)
			}
		} else {
			c.handler.HandleError(data, err)
		}
	}
}
