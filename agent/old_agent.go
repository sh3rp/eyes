package agent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/matishsiao/goInfo"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/messages"
	"github.com/sh3rp/eyes/probe"
	"github.com/sh3rp/eyes/util"
)

var VERSION_MAJOR = 0
var VERSION_MINOR = 1
var VERSION_PATCH = 0

type ProbeAgent struct {
	ID            string
	IPAddress     string
	Label         string
	Location      string
	Connection    net.Conn
	ResultChannel chan *messages.AgentProbeResult
	OOBChannel    chan *messages.AgentMessage
}

func NewAgent(label, location string) *ProbeAgent {
	return &ProbeAgent{
		ID:            util.GenID(),
		IPAddress:     util.GetLocalIP(),
		Label:         label,
		Location:      location,
		ResultChannel: make(chan *messages.AgentProbeResult),
		OOBChannel:    make(chan *messages.AgentMessage),
	}
}

func (a *ProbeAgent) connect(host string) net.Conn {
	c, err := net.Dial("tcp", host+":12121")

	if err != nil {
		log.Error().Msgf("Error connecting: %v", err)
		return nil
	}

	log.Info().Msgf("Connected: %s", host)

	return c
}

func (a *ProbeAgent) getAgentInfo() *messages.AgentInfo {
	platformInfo := goInfo.GetInfo()

	info := &messages.AgentInfo{
		Ipaddress:    a.IPAddress,
		Label:        a.Label,
		Location:     a.Location,
		AgentVersion: fmt.Sprintf("%d.%d.%d", VERSION_MAJOR, VERSION_MINOR, VERSION_PATCH),
		Hostname:     platformInfo.Hostname,
		Os:           fmt.Sprintf("%s/%s %s (%d cpus)", platformInfo.Kernel, platformInfo.Platform, platformInfo.Core, platformInfo.CPUs),
	}
	return info
}

func (a *ProbeAgent) Start(controllerHost string) {
	log.Info().Msgf("Starting agent: %s (%s) - %s", a.ID, a.Label, a.IPAddress)
	for {
		var c net.Conn

		for c == nil {
			c = a.connect(controllerHost)
			time.Sleep(5 * time.Second)
		}

		a.Connection = c

		hello := &messages.AgentMessage{
			Type: messages.AgentMessage_HELLO,
			Id:   a.ID,
			Auth: "changeit",
			Info: a.getAgentInfo(),
		}

		msg, err := proto.Marshal(hello)

		if err != nil {
			log.Error().Msgf("Error marshaling hello packet: %v", err)
			break
		}

		a.Connection.Write(msg)

		go a.WriteLoop()

		for {
			data := make([]byte, 4096)
			len, err := c.Read(data)

			if err != nil {
				log.Error().Msgf("ERROR (read): %v", err)
				break
			}

			cmd := &messages.ControllerMessage{}
			err = proto.Unmarshal(data[:len], cmd)

			if err != nil {
				log.Error().Msgf("ERROR (marshal): %v", err)
			} else {
				go a.Dispatch(cmd)
			}
		}
	}
}

func (a *ProbeAgent) WriteLoop() {
	for {
		select {
		case result := <-a.ResultChannel:
			msg := &messages.AgentMessage{
				Type:   messages.AgentMessage_RESULT,
				Result: result,
				Id:     a.ID,
			}

			data, err := proto.Marshal(msg)

			if err != nil {
				log.Error().Msgf("ERROR (writeLoop:result): %v", err)
			} else {
				a.Connection.Write(data)
			}
		case oob := <-a.OOBChannel:
			data, err := proto.Marshal(oob)

			if err != nil {
				log.Error().Msgf("ERROR (writeLoop:oob): %v", err)
			} else {
				a.Connection.Write(data)
			}
		}
	}
}

func (a *ProbeAgent) Dispatch(cmd *messages.ControllerMessage) {
	switch cmd.Type {
	case messages.ControllerMessage_AGENT_INFO_REQUEST:
		a.OOBChannel <- &messages.AgentMessage{
			Id:   a.ID,
			Info: a.getAgentInfo(),
		}
	case messages.ControllerMessage_LATENCY_REQUEST:
		req := cmd.LatencyRequest
		var port int
		if _, ok := req.Parameters["port"]; ok {
			port, _ = strconv.Atoi(req.Parameters["port"])
		} else {
			port = 80
		}
		latency, latencyErr := probe.GetLatency(a.IPAddress, req.Host, uint16(port))
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, latency)
		if err == nil {
			var success = latencyErr == nil
			errorMessage := "ok"
			if !success {
				errorMessage = latencyErr.Error()
			}
			a.ResultChannel <- &messages.AgentProbeResult{
				Host:         req.Host,
				ProbeId:      a.ID,
				Data:         buf.Bytes(),
				Timestamp:    time.Now().UnixNano(),
				Type:         messages.AgentProbeResult_TCP,
				ResultId:     req.ResultId,
				ErrorMessage: errorMessage,
				Error:        success,
			}
		} else {
			log.Error().Msgf("Error packing bytes: %v", err)
		}
	}
}
