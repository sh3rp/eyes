package controller

import (
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/messages"
)

type ProbeAgent struct {
	Id         string
	Info       *messages.AgentInfo
	Connection net.Conn
}

func (pa *ProbeAgent) SendCommand(cmd *messages.ControllerMessage) error {
	data, err := proto.Marshal(cmd)

	if err != nil {
		return err
	}

	pa.Connection.Write(data)
	return nil
}

func (pa *ProbeAgent) ReadLoop(resultChannel chan *messages.AgentProbeResult, disconnectChannel chan string) {
	for {
		data := make([]byte, 4096)
		len, err := pa.Connection.Read(data)

		if err != nil {
			log.Error().Msgf("ERROR (readLoop): %v", err)
			disconnectChannel <- pa.Id
			return
		}

		agentMessage := &messages.AgentMessage{}
		err = proto.Unmarshal(data[:len], agentMessage)

		if err != nil {
			log.Error().Msgf("ERROR (unmarshal): %v", err)
		}

		switch agentMessage.Type {
		case messages.AgentMessage_RESULT:
			resultChannel <- agentMessage.Result
		}
	}
}
