package action

import (
	"strconv"
	"strings"
	"time"

	"github.com/oklog/ulid"
	"github.com/sh3rp/eyes/agent"
)

type SSHExec struct{}

func (ssh *SSHExec) Execute(id ulid.ULID, config agent.ActionConfig) (agent.Result, error) {
	host := config.Parameters["host"]
	port := 22

	if _, ok := config.Parameters["port"]; ok {
		n, err := strconv.Atoi(config.Parameters["port"])

		if err != nil {
			return agent.Result{}, err
		}
		port = n
	}

	username := config.Parameters["username"]
	password := config.Parameters["password"]
	client := NewSSHClient(host, port, username, password)

	lines, err := client.Run(config.Parameters["command"])

	if err != nil {
		return agent.Result{}, err
	}

	str := strings.Join(lines, "\n")

	return agent.Result{
		ID:        id,
		ConfigID:  config.Id,
		Data:      []byte(str),
		Timestamp: time.Now().Unix() / 1000000,
	}, nil
}
