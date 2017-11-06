package net

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/msg"
	"github.com/sh3rp/eyes/util"
)

var DEFAULT_BUFFER_SIZE = 8128

type Connection interface {
	Send(msg.Packet)
	SetHandler(PacketHandler) error
	GetRemoteInfo() (*msg.NodeInfo, error)
}

type PacketHandler interface {
	HandlePacket(msg.Packet)
	HandleError([]byte, error)
}

type connection struct {
	conn                net.Conn
	bufferSize          int
	handler             PacketHandler
	localInfo           *msg.NodeInfo
	remoteInfo          *msg.NodeInfo
	packetQueue         chan msg.Packet
	lastKeepalive       int64
	keepAliveTimeout    int64
	numFailedKeepalives int
}

func NewConnection(conn net.Conn, info msg.NodeInfo) Connection {
	c := &connection{
		conn:                conn,
		bufferSize:          DEFAULT_BUFFER_SIZE,
		localInfo:           &info,
		packetQueue:         make(chan msg.Packet, 10),
		keepAliveTimeout:    10000, // 10 seconds
		numFailedKeepalives: 3,
		lastKeepalive:       util.Now(),
	}

	return c
}

func (c *connection) SetHandler(handler PacketHandler) error {
	c.handler = handler
	go c.read()
	go c.write()
	//go c.keepalive()
	return nil
}

func (c *connection) Send(pkt msg.Packet) {
	c.packetQueue <- pkt
}

func (c *connection) GetRemoteInfo() (*msg.NodeInfo, error) {
	return c.remoteInfo, nil
}

func (c *connection) write() {
	for pkt := range c.packetQueue {
		if c.conn == nil {
			errors.New("Connection not set, cannot write packet")
			continue
		}
		data, err := proto.Marshal(&pkt)
		if err != nil {
			log.Error().Msgf("Error marshalling: %v", err)
		}
		_, err = c.conn.Write(data)
		if err != nil {
			log.Error().Msgf("Error writing: %v", err)
		}
	}
}

func (c *connection) read() {
	if c.conn == nil {
		c.handler.HandleError(nil, errors.New("No connection set, cannot read"))
		return
	}
	if c.handler == nil {
		log.Error().Msgf("Catastrophic error: No handler set for inbound data")
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
				if pkt.Type == msg.Packet_KEEPALIVE {
					c.remoteInfo = pkt.GetInfo()
					c.lastKeepalive = util.Now()
					log.Info().Msgf("Keepalive: id=%s v%d.%d.%d",
						c.remoteInfo.Id,
						c.remoteInfo.MajorVersion,
						c.remoteInfo.MinorVersion,
						c.remoteInfo.PatchVersion)
				} else {
					c.handler.HandlePacket(pkt)
				}
			}
		} else {
			c.handler.HandleError(data, err)
		}
	}
}

func (c *connection) keepalive() {
	for {
		if (c.lastKeepalive + c.keepAliveTimeout) <= util.Now() {
			fmt.Printf("ka1\n")
			c.Send(msg.Packet{
				Type:   msg.Packet_KEEPALIVE,
				Packet: &msg.Packet_Info{c.localInfo},
			})
			fmt.Printf("ka2\n")
		}
		time.Sleep(time.Duration(c.keepAliveTimeout+int64(rand.Int()%5000)) * time.Millisecond)
	}
}
