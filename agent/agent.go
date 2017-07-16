package agent

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/oklog/ulid"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/messages"
	"github.com/sh3rp/eyes/probe"
)

type ProbeAgent struct {
	ID            string
	Label         string
	Connection    net.Conn
	ResultChannel chan *messages.ProbeResult
}

func NewAgent(label string) *ProbeAgent {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return &ProbeAgent{
		ID:            id.String(),
		Label:         label,
		ResultChannel: make(chan *messages.ProbeResult),
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

func (a *ProbeAgent) Start(controllerHost string) {
	log.Info().Msgf("Starting agent: %s (%s)", a.ID, a.Label)
	for {
		var c net.Conn

		for c == nil {
			c = a.connect(controllerHost)
			time.Sleep(5 * time.Second)
		}

		a.Connection = c

		hello := &messages.ProbeACK{
			Type:  messages.ProbeACK_HELLO,
			Id:    a.ID,
			Label: a.Label,
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

			cmd := &messages.ProbeCommand{}
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
		result := <-a.ResultChannel

		msg := &messages.ProbeACK{
			Type:   messages.ProbeACK_RESULT,
			Result: result,
			Id:     a.ID,
		}

		data, err := proto.Marshal(msg)

		if err != nil {
			log.Error().Msgf("ERROR (writeLoop): %v", err)
		} else {
			a.Connection.Write(data)
		}
	}
}

func (a *ProbeAgent) Dispatch(cmd *messages.ProbeCommand) {
	switch cmd.Type {
	// used for testing purposes
	case messages.ProbeCommand_NOOP:
		a.ResultChannel <- &messages.ProbeResult{
			ProbeId:   a.ID,
			Data:      []byte{0, 1, 2, 3, 4, 5, 6, 7},
			Type:      messages.ProbeResult_NOOP,
			Timestamp: time.Now().UnixNano(),
			Host:      GetLocalIP(),
		}
	// run TCP probe
	case messages.ProbeCommand_TCP:
		var port int
		if _, ok := cmd.Parameters["port"]; ok {
			port, _ = strconv.Atoi(cmd.Parameters["port"])
		} else {
			port = 80
		}
		latency := probe.GetLatency(GetLocalIP(), cmd.Host, uint16(port))
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, latency)
		if err == nil {
			a.ResultChannel <- &messages.ProbeResult{
				Host:      cmd.Host,
				ProbeId:   a.ID,
				Data:      buf.Bytes(),
				Timestamp: time.Now().UnixNano(),
				Type:      messages.ProbeResult_TCP,
			}
		} else {
			log.Error().Msgf("Error packing bytes: %v", err)
		}
	}
}

func GetLocalIP() string {
	var ip string
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if addr.String() != "127.0.0.1" && !strings.Contains(addr.String(), ":") {
			ipAddr := addr.String()
			elements := strings.Split(ipAddr, "/")
			ip = elements[0]
		}
	}
	return ip
}
