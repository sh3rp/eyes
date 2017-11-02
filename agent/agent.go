package agent

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	cron "gopkg.in/robfig/cron.v2"
)

type Agent interface {
	StoreActionConfig(ActionConfig) error
	DeleteActionConfig(string) error
	ScheduleAction(string, string) error
	UnscheduleAction(string) error
}

type MemAgent struct {
	Configs       map[string]ActionConfig
	ResultFuncs   []func(Result)
	schedules     map[string]cron.EntryID
	cronScheduler *cron.Cron
	results       chan Result
}

func NewMemAgent() *MemAgent {
	cronScheduler := cron.New()
	cronScheduler.Start()
	agent := &MemAgent{
		Configs:       make(map[string]ActionConfig),
		schedules:     make(map[string]cron.EntryID),
		cronScheduler: cronScheduler,
		results:       make(chan Result),
	}
	go agent.sendResults()
	return agent
}

func (a *MemAgent) StoreActionConfig(c ActionConfig) error {
	a.Configs[c.Id] = c
	return nil
}

func (a *MemAgent) DeleteActionConfig(id string) error {
	delete(a.Configs, id)
	return nil
}

func (a *MemAgent) ScheduleAction(id, cronString string) error {
	config := a.Configs[id]
	action := ACTIONS[config.Action]
	scheduleId, err := a.cronScheduler.AddFunc(cronString, func() {
		result, err := action.Execute(config)
		if err != nil {
			log.Error().Msgf("Error executing config %s", id)
		} else {
			a.results <- result
		}
	})
	if err != nil {
		log.Error().Msgf("Error scheduling: %v", err)
	} else {
		a.schedules[id] = scheduleId
	}
	return err
}

func (a *MemAgent) UnscheduleAction(id string) error {
	if _, ok := a.schedules[id]; !ok {
		return errors.New(fmt.Sprintf("No schedule with ID %s", id))
	}
	a.cronScheduler.Remove(a.schedules[id])
	return nil
}

func (a *MemAgent) sendResults() {
	for result := range a.results {
		for _, f := range a.ResultFuncs {
			f(result)
		}
	}
}
