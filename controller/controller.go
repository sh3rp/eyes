package controller

import (
	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/db"
	"github.com/sh3rp/eyes/net"
)

type Controller struct {
	DB       db.EyesDB
	Listener net.Connection
	Agents   map[string]agent.Agent
}

func (ctrl Controller) NewConfig(c db.Config) error {
	return nil
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
