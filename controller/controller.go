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
	go ctrl.manageDeployments()
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
	if d.Id == "" {
		d.Id = util.NewId()
	}

	err := ctrl.database.SaveDeployment(d)

	if err != nil {
		return err
	}

	return nil
}

func (ctrl Controller) GetDeployments() ([]db.Deployment, error) {
	return ctrl.database.GetDeployments()
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

func (ctrl Controller) GetAgent(id util.ID) (db.Agent, error) {
	return ctrl.database.GetAgent(id)
}

func (ctrl Controller) manageAgents() {
	for {
		persistedAgents, _ := ctrl.database.GetAgents()

		for _, a := range persistedAgents {
			if _, ok := ctrl.agents[a.Id]; !ok {
				if a.AgentType == 0 {
					fmt.Printf("Deploying agent: %s\n", a.Id)
					ctrl.agents[a.Id] = agent.NewMemAgent(ctrl.resultHandler)
				}
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (ctrl Controller) manageDeployments() {
	for {
		for id, _ := range ctrl.agents {
			ctrl.manageAgentDeployment(id)
		}
		time.Sleep(5000 * time.Millisecond)
	}
}

func (ctrl Controller) manageAgentDeployment(agentId util.ID) {
	var thisAgent agent.Agent
	var ok bool
	if thisAgent, ok = ctrl.agents[agentId]; !ok {
		fmt.Printf("Error: no such agent configured")
		return
	}

	deployments, err := ctrl.database.GetDeployments()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var agentSchedules []db.Schedule

	for _, deployment := range deployments {
		if deployment.Agent == agentId {
			schedule, err := ctrl.database.GetSchedule(deployment.Schedule)
			if err == nil {
				agentSchedules = append(agentSchedules, schedule)
			}
		}
	}

	for _, agentSchedule := range agentSchedules {
		config, err := ctrl.database.GetConfig(agentSchedule.ConfigId)

		if err != nil {
			fmt.Printf("Error: no such config: %v\n", err)
		}

		thisAgent.StoreConfig(config)
		thisAgent.ScheduleConfig(config.Id, agentSchedule.Schedule)
	}

}
