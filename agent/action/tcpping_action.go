package action

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/oklog/ulid"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/probe"
)

type TCPPing struct{}

func (ping *TCPPing) Execute(id ulid.ULID, config agent.ActionConfig) (agent.Result, error) {
	var result agent.Result
	var port int
	if _, ok := config.Parameters["port"]; ok {
		port, _ = strconv.Atoi(config.Parameters["port"])
	} else {
		port = 80
	}
	latency, latencyErr := probe.GetLatency(config.Parameters["srcIp"], config.Parameters["dstIp"], uint16(port))
	if latencyErr != nil {
		return agent.Result{
			ID:        id,
			Data:      []byte{},
			DataCode:  agent.DATA_ERROR,
			Timestamp: agent.Now(),
		}, nil
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, latency)
	if err == nil {
		result = agent.Result{
			ID:        id,
			Data:      buf.Bytes(),
			DataCode:  agent.DATA_OK,
			Timestamp: agent.Now(),
		}
	} else {
		log.Error().Msgf("Error packing bytes: %v", err)
		result = agent.Result{
			ID:        id,
			Data:      []byte{},
			DataCode:  agent.DATA_ERROR,
			Timestamp: agent.Now(),
		}
	}
	return result, nil
}
