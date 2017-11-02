package agent

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/probe"
)

type TCPPing struct{}

func (ping TCPPing) Execute(config ActionConfig) (Result, error) {
	var result Result
	var port int
	if _, ok := config.Parameters["port"]; ok {
		port, _ = strconv.Atoi(config.Parameters["port"])
	} else {
		port = 80
	}
	latency, latencyErr := probe.GetLatency(config.Parameters["srcIp"], config.Parameters["dstIp"], uint16(port))
	if latencyErr != nil {
		return Result{
			Id:        NewId(),
			Data:      []byte{},
			DataCode:  DATA_ERROR,
			Timestamp: Now(),
		}, nil
	}
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, latency)
	if err == nil {
		result = Result{
			Id:        NewId(),
			Data:      buf.Bytes(),
			DataCode:  DATA_OK,
			Timestamp: Now(),
		}
	} else {
		log.Error().Msgf("Error packing bytes: %v", err)
		result = Result{
			Id:        NewId(),
			Data:      []byte{},
			DataCode:  DATA_ERROR,
			Timestamp: Now(),
		}
	}
	return result, nil
}
