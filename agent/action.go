package agent

import (
	"github.com/sh3rp/eyes/db"
	"github.com/sh3rp/eyes/util"
)

type Action interface {
	Execute(db.Config) (Result, error)
}

const (
	DATA_OK = iota
	DATA_WARN
	DATA_ERROR
)

const (
	A_TEST = iota
	A_TCPPING
	A_SSH
	A_SNMP
	A_HTTP
)

var ACTIONS = map[int]Action{
	A_TEST:    DummyAction{},
	A_TCPPING: TCPPing{},
	A_SSH:     SSHExec{},
	A_SNMP:    SNMPPoll{},
	A_HTTP:    WebAction{},
}

type Result struct {
	Id        util.ID
	ConfigId  util.ID
	DataCode  int
	Data      []byte
	Tags      map[string]string
	Timestamp int64
}

type DummyAction struct{}

func (d DummyAction) Execute(c db.Config) (Result, error) {
	return Result{
		Id:        util.NewId(),
		ConfigId:  c.Id,
		DataCode:  DATA_OK,
		Data:      []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		Tags:      c.Parameters,
		Timestamp: util.Now(),
	}, nil
}
