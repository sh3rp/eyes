package controller

import (
	"github.com/sh3rp/eyes/agent"
)

type ControllerServer interface {
	SayHello(Auth)
	ProcessResult(agent.Result)
}

type Auth struct {
	Username string
	Password string
}
