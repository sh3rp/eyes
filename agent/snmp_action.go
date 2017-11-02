package agent

import "github.com/rs/zerolog/log"

type SNMPPoll struct{}

func (snmp SNMPPoll) Execute(config ActionConfig) (Result, error) {
	log.Info().Msgf("Polling snmp")
	return Result{}, nil
}
