package controller

import (
	"github.com/jasonlvhit/gocron"
	"github.com/sh3rp/eyes/messages"
)

type AgentScheduler struct {
	Schedulers map[string]*gocron.Scheduler
	Controller *ProbeController
}

func NewAgentScheduler(controller *ProbeController) *AgentScheduler {
	defaultScheduler := gocron.NewScheduler()
	schedulers := make(map[string]*gocron.Scheduler)
	schedulers["default"] = defaultScheduler
	return &AgentScheduler{
		Schedulers: schedulers,
		Controller: controller,
	}
}

func (a *AgentScheduler) ScheduleEveryXSeconds(numSeconds uint64, agentId string, cmd *messages.ProbeCommand) {
	a.Schedulers["default"].Every(numSeconds).Seconds().Do(a.Controller.SendProbe, agentId, cmd)
}
