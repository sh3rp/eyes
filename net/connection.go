package net

import (
	"net"

	"github.com/gogo/protobuf/proto"
	"github.com/sh3rp/eyes/msg"
)

var DEFAULT_BUFFER_SIZE = 8128

type Connection interface {
	Send(msg.Packet) error
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

func NewConnection(conn net.Conn, handler PacketHandler, readBufSize int) Connection {
	bufSize := DEFAULT_BUFFER_SIZE
	if readBufSize > 0 {
		bufSize = readBufSize
	}
	c := &connection{
		conn:       conn,
		handler:    handler,
		bufferSize: bufSize,
	}
	go c.read()

	return c
}

func (c *connection) Send(pkt msg.Packet) error {
	return c.write(pkt)
}

func (c *connection) write(pkt msg.Packet) error {
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
	if c.conn == nil {
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
