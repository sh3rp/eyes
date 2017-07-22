package controller

import (
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/messages"
)

type ProbeAgent struct {
	Id         string
	IPAddress  string
	Label      string
	Location   string
	Connection net.Conn
}

func (pa *ProbeAgent) SendCommand(cmd *messages.ProbeCommand) error {
	data, err := proto.Marshal(cmd)

	if err != nil {
		return err
	}

	pa.Connection.Write(data)
	return nil
}

func (pa *ProbeAgent) ReadLoop(resultChannel chan *messages.ProbeResult, disconnectChannel chan string) {
	for {
		data := make([]byte, 4096)
		len, err := pa.Connection.Read(data)

		if err != nil {
			log.Error().Msgf("ERROR (readLoop): %v", err)
			disconnectChannel <- pa.Id
			return
		}

		ack := &messages.ProbeACK{}
		err = proto.Unmarshal(data[:len], ack)

		if err != nil {
			log.Error().Msgf("ERROR (unmarshal): %v", err)
		}

		switch ack.Type {
		case messages.ProbeACK_RESULT:
			resultChannel <- ack.Result
		}
	}
}
