package agent

import (
	"strconv"
	"strings"

	"github.com/sh3rp/eyes/util"
)

type SSHExec struct{}

func (ssh SSHExec) Execute(config ActionConfig) (Result, error) {
	host := config.Parameters["host"]
	port := 22

	if _, ok := config.Parameters["port"]; ok {
		n, err := strconv.Atoi(config.Parameters["port"])

		if err != nil {
			return Result{}, err
		}
		port = n
	}

	username := config.Parameters["username"]
	password := config.Parameters["password"]
	client := NewSSHClient(host, port, username, password)

	lines, err := client.Run(config.Parameters["command"])

	if err != nil {
		return Result{}, err
	}

	str := strings.Join(lines, "\n")

	return Result{
		Id:        util.NewId(),
		ConfigId:  config.Id,
		Data:      []byte(str),
		Timestamp: util.Now(),
	}, nil
}
