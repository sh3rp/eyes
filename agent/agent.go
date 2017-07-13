package client

import (
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/sh3rp/eyes/messages"
	"github.com/twinj/uuid"
)

type ProbeAgent struct {
	ID            string
	Connection    net.Conn
	ResultChannel chan *messages.ProbeResult
}

func NewAgent() *ProbeAgent {
	return &ProbeAgent{
		ID:            uuid.NewV4().String(),
		ResultChannel: make(chan *messages.ProbeResult),
	}
}

func (a *ProbeAgent) Start(controllerHost string) {
	c, err := net.Dial("tcp", controllerHost+":12121")

	if err != nil {
		log.Fatal(err)
	}

	a.Connection = c

	hello := &messages.ProbeACK{
		Type: messages.ProbeACK_HELLO,
		Id:   a.ID,
	}

	msg, err := proto.Marshal(hello)

	a.Connection.Write(msg)

	for {
		data := make([]byte, 4096)
		len, err := c.Read(data)

		if err != nil {
			log.Printf("ERROR (read): %v", err)
			return
		}

		cmd := &messages.ProbeCommand{}
		err = proto.Unmarshal(data[:len], cmd)

		if err != nil {
			log.Printf("ERROR (marshal): %v", err)
		} else {
			go a.Dispatch(cmd)
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
			log.Printf("ERROR (writeLoop): %v", err)
		} else {
			a.Connection.Write(data)
		}
	}
}

func (a *ProbeAgent) Dispatch(cmd *messages.ProbeCommand) {
	switch cmd.Type {
	case messages.ProbeCommand_TCP:
		log.Printf("Sending TCP ping")
		a.ResultChannel <- &messages.ProbeResult{
			Host:      "127.0.0.1",
			ProbeId:   a.ID,
			Datapoint: 23,
		}
	}
}
