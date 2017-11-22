package agent

import (
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/db"
)

type SNMPPoll struct{}

func (snmp SNMPPoll) Execute(config db.Config) (Result, error) {
	log.Info().Msgf("Polling snmp")
	return Result{}, nil
}
