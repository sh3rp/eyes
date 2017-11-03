package agent

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid"
)

type Action interface {
	Execute(ActionConfig) (Result, error)
}

type ActionConfig struct {
	Id         string
	Action     int
	Parameters map[string]string
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
	Id        string
	ConfigId  string
	DataCode  int
	Data      []byte
	Tags      map[string]string
	Timestamp int64
}

func Now() int64 {
	return time.Now().UnixNano() / 1000000
}

func NewId() string {
	id := ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader)
	return id.String()
}

type DummyAction struct{}

func (d DummyAction) Execute(c ActionConfig) (Result, error) {
	return Result{
		Id:        NewId(),
		ConfigId:  c.Id,
		DataCode:  DATA_OK,
		Data:      []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		Tags:      c.Parameters,
		Timestamp: Now(),
	}, nil
}
