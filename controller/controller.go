package controller

import (
	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/db"
	"github.com/sh3rp/eyes/net"
	"github.com/sh3rp/eyes/util"
)

type Controller struct {
	DB       db.EyesDB
	Listener net.Connection
	Agents   map[string]agent.Agent
}

func NewController(dir string) Controller {
	bolt := db.NewDB(dir, "eyes")
	db := db.NewBoltEyesDB(bolt)
	return Controller{
		DB:     db,
		Agents: make(map[string]agent.Agent),
	}
}

func (ctrl Controller) NewConfig(c db.Config) (db.Config, error) {
	if c.Id == "" {
		c.Id = util.NewId()
	}
	err := ctrl.DB.SaveConfig(c)

	if err != nil {
		return db.Config{}, err
	}

	return ctrl.DB.GetConfig(c.Id)
}

func (ctrl Controller) GetConfigs() ([]db.Config, error) {
	return ctrl.DB.GetConfigs()
}

func (ctrl Controller) NewSchedule(s db.Schedule) error {
	return nil
}

func (ctrl Controller) NewDeployment(d db.Deployment) error {
	return nil
}

func (ctrl Controller) NewAgent(a db.Agent) error {
	return ctrl.DB.SaveAgent(a)
}

func (ctrl Controller) processAgents() {

}
