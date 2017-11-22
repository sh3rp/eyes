package agent

import (
	"github.com/sh3rp/eyes/db"
	"github.com/sh3rp/eyes/util"
)

type AgentType int

const (
	AT_LOCAL = iota
	AT_REMOTE
)

type Agent interface {
	GetType() AgentType

	StoreConfig(db.Config) error
	DeleteConfig(util.ID) error
	GetAllConfigs() ([]db.Config, error)
	ScheduleConfig(util.ID, string) error
	UnScheduleConfig(util.ID) error
	GetAllSchedules() ([]db.Schedule, error)
	HandleResult(func(Result))
}
