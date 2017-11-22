package agent

import (
	"github.com/sh3rp/eyes/db"
	"github.com/sh3rp/eyes/util"
)

type Agent interface {
	StoreConfig(db.Config) error
	DeleteConfig(util.ID) error
	GetAllConfigs() ([]db.Config, error)
	ScheduleConfig(util.ID, string) error
	UnScheduleConfig(util.ID) error
	GetAllSchedules() ([]db.Schedule, error)
	HandleResult(func(Result))
}
