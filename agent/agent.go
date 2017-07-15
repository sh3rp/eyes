package agent

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/oklog/ulid"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/messages"
	"github.com/sh3rp/eyes/probe"
)

type ProbeAgent struct {
	ID            string
	Connection    net.Conn
	ResultChannel chan *messages.ProbeResult
}

func NewAgent() *ProbeAgent {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return &ProbeAgent{
		ID:            id.String(),
		ResultChannel: make(chan *messages.ProbeResult),
	}
}

func (a *ProbeAgent) connect(host string) net.Conn {
	c, err := net.Dial("tcp", host+":12121")

	if err != nil {
		log.Debug().Msgf("Error connecting: %v", err)
		return nil
	}

	log.Info().Msgf("Connected: %s", host)

	return c
}

func (a *ProbeAgent) Start(controllerHost string) {
	for {
		var c net.Conn

		for c == nil {
			c = a.connect(controllerHost)
			time.Sleep(5 * time.Second)
		}

		a.Connection = c

		hello := &messages.ProbeACK{
			Type: messages.ProbeACK_HELLO,
			Id:   a.ID,
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
		}

		data, err := proto.Marshal(msg)

		if err != nil {
			log.Error().Msgf("ERROR (writeLoop): %v", err)
		} else {
			log.Debug().Msgf("Sending result back")
			a.Connection.Write(data)
		}
	}
}

func (a *ProbeAgent) Dispatch(cmd *messages.ProbeCommand) {
	switch cmd.Type {
	case messages.ProbeCommand_NOOP:
		a.ResultChannel <- &messages.ProbeResult{
			ProbeId:   a.ID,
			Data:      []byte{0, 1, 2, 3, 4, 5, 6, 7},
			Type:      messages.ProbeResult_NOOP,
			Timestamp: time.Now().UnixNano(),
			Host:      "127.0.0.1",
		}
	case messages.ProbeCommand_TCP:
		latency := probe.GetLatency("127.0.0.1", cmd.Host, 80)
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, latency)
		if err == nil {
			a.ResultChannel <- &messages.ProbeResult{
				Host:    "127.0.0.1",
				ProbeId: a.ID,
				Data:    buf.Bytes(),
			}
		} else {
			log.Error().Msgf("Error packing bytes: %v", err)
		}
	}
}
