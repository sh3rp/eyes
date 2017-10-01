package agent

import (
	"time"

	"github.com/oklog/ulid"
)

type Probe map[Action]Result

type Action interface {
	Execute(ulid.ULID, ActionConfig) (Result, error)
}

type ActionConfig struct {
	Id         ulid.ULID
	Parameters map[string]string
}

const (
	DATA_OK = iota
	DATA_WARN
	DATA_ERROR
)

type Result struct {
	ID        ulid.ULID
	ConfigID  ulid.ULID
	DataCode  int
	Data      []byte
	Timestamp int64
}

func Now() int64 {
	return time.Now().UnixNano() / 1000000
}
