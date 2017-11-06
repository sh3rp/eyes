package agent

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/util"
	cron "gopkg.in/robfig/cron.v2"
)

type Agent interface {
	StoreActionConfig(ActionConfig) error
	DeleteActionConfig(util.ID) error
	ScheduleAction(util.ID, string) error
	UnscheduleAction(util.ID) error
	AddResultHandler(string, func(Result)) error
	RemoveResultHandler(string) error
}

type MemAgent struct {
	Configs       map[util.ID]ActionConfig
	resultFuncs   map[string]func(Result)
	schedules     map[util.ID]cron.EntryID
	cronScheduler *cron.Cron
	results       chan Result
}

func NewMemAgent() *MemAgent {
	cronScheduler := cron.New()
	cronScheduler.Start()
	agent := &MemAgent{
		Configs:       make(map[util.ID]ActionConfig),
		schedules:     make(map[util.ID]cron.EntryID),
		cronScheduler: cronScheduler,
		results:       make(chan Result),
		resultFuncs:   make(map[string]func(Result)),
	}
	go agent.sendResults()
	return agent
}

func (a *MemAgent) StoreActionConfig(c ActionConfig) error {
	a.Configs[c.Id] = c
	return nil
}

func (a *MemAgent) DeleteActionConfig(id util.ID) error {
	delete(a.Configs, id)
	return nil
}

func (a *MemAgent) ScheduleAction(id util.ID, cronString string) error {
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

func (a *MemAgent) UnscheduleAction(id util.ID) error {
	if _, ok := a.schedules[id]; !ok {
		return errors.New(fmt.Sprintf("No schedule with ID %s", id))
	}
	a.cronScheduler.Remove(a.schedules[id])
	return nil
}

func (a *MemAgent) AddResultHandler(id string, f func(Result)) error {
	a.resultFuncs[id] = f
	return nil
}

func (a *MemAgent) RemoveResultHandler(id string) error {
	delete(a.resultFuncs, id)
	return nil
}

func (a *MemAgent) sendResults() {
	for result := range a.results {
		for _, f := range a.resultFuncs {
			f(result)
		}
	}
}
