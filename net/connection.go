package net

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/msg"
	"github.com/sh3rp/eyes/util"
)

var D = func(marker string, v interface{}) { fmt.Printf("%s: %v\n", marker, v) }

var DEFAULT_BUFFER_SIZE = 8128

type Connection interface {
	Send(msg.Packet)
	SetHandler(PacketHandler) error
	GetRemoteInfo() (msg.NodeInfo, error)
}

type PacketHandler interface {
	HandlePacket(msg.Packet)
	HandleError([]byte, error)

	HandleHello(msg.Hello)
	HandleKeepalive(msg.KeepAlive)
	HandleScheduleActionConfig(msg.ScheduleActionConfig)
	HandleResult(msg.Result)
	HandleUnScheduleActionConfig(msg.UnscheduleActionConfig)
	HandleRunActionConfig(msg.RunActionConfig)
	HandleAllActionConfigs(msg.AllActionConfigs)
}

type connection struct {
	conn                net.Conn
	bufferSize          int
	handler             PacketHandler
	localInfo           msg.NodeInfo
	remoteInfo          msg.NodeInfo
	packetQueue         chan msg.Packet
	lastKeepalive       int64
	keepAliveTimeout    int64
	numFailedKeepalives int
}

func NewConnection(conn net.Conn, info msg.NodeInfo, keepaliveTimeout int64) Connection {
	c := &connection{
		conn:                conn,
		bufferSize:          DEFAULT_BUFFER_SIZE,
		localInfo:           info,
		packetQueue:         make(chan msg.Packet, 10),
		keepAliveTimeout:    keepaliveTimeout, // 10 seconds
		numFailedKeepalives: 3,
		lastKeepalive:       util.Now(),
	}

	return c
}

func (c *connection) SetHandler(handler PacketHandler) error {
	c.handler = handler
	go c.read()
	go c.write()
	go c.keepalive()
	return nil
}

func (c *connection) Send(pkt msg.Packet) {
	c.packetQueue <- pkt
}

func (c *connection) GetRemoteInfo() (msg.NodeInfo, error) {
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
				switch pkt.Packet.(type) {
				case *msg.Packet_Keepalive:
					c.remoteInfo = *pkt.GetKeepalive().Info

					c.lastKeepalive = util.Now()

					log.Debug().Msgf("Keepalive: id=%s v%d.%d.%d (%s %s/%s (%s))",
						c.remoteInfo.Id,
						c.remoteInfo.MajorVersion,
						c.remoteInfo.MinorVersion,
						c.remoteInfo.PatchVersion,
						c.remoteInfo.Hostname,
						c.remoteInfo.Os,
						c.remoteInfo.Kernel,
						c.remoteInfo.Platform)
					c.handler.HandleKeepalive(*pkt.GetKeepalive())
				case *msg.Packet_Hello:
					c.handler.HandleHello(*pkt.GetHello())
				case *msg.Packet_Schedule:
					c.handler.HandleScheduleActionConfig(*pkt.GetSchedule())
				case *msg.Packet_Unschedule:
					c.handler.HandleUnScheduleActionConfig(*pkt.GetUnschedule())
				case *msg.Packet_Run:
					c.handler.HandleRunActionConfig(*pkt.GetRun())
				case *msg.Packet_Result:
					c.handler.HandleResult(*pkt.GetResult())
				case *msg.Packet_AllConfigs:
					c.handler.HandleAllActionConfigs(*pkt.GetAllConfigs())
				case nil:
					log.Error().Msg("ERROR: empty packet (nil)")
				default:
					log.Error().Msgf("ERROR: unknown packet: %v", pkt)
				}

				c.handler.HandlePacket(pkt)
			}
		} else {
			c.handler.HandleError(data, err)
		}
	}
}

func (c *connection) keepalive() {
	for {
		if (c.lastKeepalive + c.keepAliveTimeout) <= util.Now() {
			p := msg.Packet{
				Packet: &msg.Packet_Keepalive{
					&msg.KeepAlive{
						Info:      &c.localInfo,
						Timestamp: util.Now(),
					},
				},
				Code: 0,
				Msg:  "ok",
			}
			c.Send(p)
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}
