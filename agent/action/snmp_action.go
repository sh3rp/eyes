package action

import (
	"github.com/oklog/ulid"
	"github.com/sh3rp/eyes/agent"
)

type SNMPPoll struct{}

func (snmp *SNMPPoll) Execute(id ulid.ULID, config agent.ActionConfig) (agent.Result, error) {
	return agent.Result{}, nil
}
