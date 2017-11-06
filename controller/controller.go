package controller

import (
	"github.com/sh3rp/eyes/agent"
)

type Controller interface {
	SayHello(Auth)
	ProcessResult(agent.Result)
}

type Auth struct {
	Username string
	Password string
}
