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

func (c Controller) NewConfig(c db.Config) error {

}

func (c Controller) NewSchedule(s db.Schedule) error {

}

func (c Controller) NewDeployment(d db.Deployment) error {

}

func (c Controller) NewAgent(a db.Agent) error {
	return c.DB.SaveAgent(a)
}

func (c Controller) processAgents() {

}
