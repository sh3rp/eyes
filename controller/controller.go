package controller

import (
	"fmt"
	"time"

	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/db"
	"github.com/sh3rp/eyes/net"
	"github.com/sh3rp/eyes/util"
)

type Controller struct {
	database      db.EyesDB
	listener      net.Connection
	agents        map[util.ID]agent.Agent
	resultHandler func(agent.Result)
}

func NewController(dir string, resultHandler func(agent.Result)) Controller {
	bolt := db.NewDB(dir, "eyes")
	db := db.NewBoltEyesDB(bolt)
	ctrl := Controller{
		database:      db,
		agents:        make(map[util.ID]agent.Agent),
		resultHandler: resultHandler,
	}
	go ctrl.manageAgents()
	return ctrl
}

func (ctrl Controller) NewConfig(c db.Config) (db.Config, error) {
	if c.Id == "" {
		c.Id = util.NewId()
	}
	err := ctrl.database.SaveConfig(c)

	if err != nil {
		return db.Config{}, err
	}

	return ctrl.database.GetConfig(c.Id)
}

func (ctrl Controller) GetConfigs() ([]db.Config, error) {
	return ctrl.database.GetConfigs()
}

func (ctrl Controller) NewSchedule(s db.Schedule) (db.Schedule, error) {
	if s.Id == "" {
		s.Id = util.NewId()
	}

	err := ctrl.database.SaveSchedule(s)

	if err != nil {
		return db.Schedule{}, err
	}

	return ctrl.database.GetSchedule(s.Id)
}

func (ctrl Controller) GetSchedules() ([]db.Schedule, error) {
	return ctrl.database.GetSchedules()
}

func (ctrl Controller) NewDeployment(d db.Deployment) error {
	return nil
}

func (ctrl Controller) NewAgentLocal() (db.Agent, error) {
	a := db.Agent{Id: util.NewId(), AgentType: 0}

	err := ctrl.database.SaveAgent(a)

	if err != nil {
		return db.Agent{}, err
	}

	return ctrl.database.GetAgent(a.Id)
}

func (ctrl Controller) GetAgents() ([]db.Agent, error) {
	return ctrl.database.GetAgents()
}

func (ctrl Controller) manageAgents() {
	for {
		persistedAgents, _ := ctrl.database.GetAgents()

		for _, a := range persistedAgents {
			if _, ok := ctrl.agents[a.Id]; !ok {
				if a.AgentType == 0 {
					ctrl.agents[a.Id] = agent.NewMemAgent(ctrl.resultHandler)
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (ctrl Controller) manageDeployments() {
	for {
		deployments, _ := ctrl.database.GetDeployments()

		schedules := make(map[util.ID][]db.Schedule)

		for id, a := range ctrl.agents {
			scheds, _ := a.GetAllSchedules()
			schedules[id] = scheds
		}

		for _, d := range deployments {
			fmt.Println("%+v", d)
		}
	}
}
