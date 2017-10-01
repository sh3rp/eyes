package agent

import (
	"github.com/oklog/ulid"
	"github.com/sh3rp/eyes/msg"
)

type ProbeListener struct {
	ID             ulid.ULID
	ListenerAction func(msg.ProbeResult) error
}

type Agent interface {
	SayHello(msg.Hello) error
	ConfigureProbe(msg.ProbeConfig) error
	ActivateProbe(ulid.ULID) error
	DeactivateProbe(ulid.ULID) error
	DeleteProbe(ulid.ULID) error
	RunProbeOnce(ulid.ULID) error
}

type ConfigEntry struct {
	Id      ulid.ULID
	Config  msg.ProbeConfig
	Active  bool
	LastRun int64
}
