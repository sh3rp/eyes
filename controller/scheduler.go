package controller

import (
	"github.com/chacken/gocron"
	"github.com/sh3rp/eyes/messages"
)

type AgentScheduler struct {
	Schedulers map[string]*gocron.Scheduler
	Controller *ProbeController
}

func NewAgentScheduler(controller *ProbeController) *AgentScheduler {
	schedulers := make(map[string]*gocron.Scheduler)
	schedulers["default"] = gocron.NewScheduler()
	go schedulers["default"].Start()
	return &AgentScheduler{
		Schedulers: schedulers,
		Controller: controller,
	}
}

func (a *AgentScheduler) ScheduleEveryXSeconds(numSeconds uint64, agentId string, cmd *messages.ControllerLatencyRequest) {
	a.Schedulers["default"].Job(cmd.ResultId).Every(1).Second().Do(a.Controller.SendProbe, agentId, cmd)
}

func (a *AgentScheduler) Cancel(cmdId string) {
	a.Schedulers["default"].Remove(cmdId)
}
