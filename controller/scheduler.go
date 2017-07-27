package controller

import (
	"github.com/chacken/gocron"
	"github.com/sh3rp/eyes/messages"
)

type AgentScheduler struct {
	Schedulers map[string]*gocron.Scheduler
	Controller *ProbeController
}

type AgentJob struct {
	AgentId    string
	Command    *messages.ProbeCommand
	Controller *ProbeController
}

func NewAgentJob(controller *ProbeController, agentId string, cmd *messages.ProbeCommand) *AgentJob {
	return &AgentJob{
		AgentId:    agentId,
		Command:    cmd,
		Controller: controller,
	}
}

func (aj *AgentJob) Run() {
	aj.Controller.SendProbe(aj.AgentId, aj.Command)
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

func (a *AgentScheduler) ScheduleEveryXSeconds(numSeconds uint64, agentId string, cmd *messages.ProbeCommand) {
	job := NewAgentJob(a.Controller, agentId, cmd)
	a.Schedulers["default"].Job(cmd.Id).Every(1).Second().Do(job.Run)
}

func (a *AgentScheduler) Cancel(cmdId string) {
	a.Schedulers["default"].Remove(cmdId)
}
