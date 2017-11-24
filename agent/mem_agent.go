package agent

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/db"
	"github.com/sh3rp/eyes/util"
	cron "gopkg.in/robfig/cron.v2"
)

type MemAgent struct {
	Configs         map[util.ID]db.Config
	resultHandler   func(Result)
	schedules       map[util.ID]cron.EntryID
	scheduleStrings map[cron.EntryID]string
	cronScheduler   *cron.Cron
	results         chan Result
}

func NewMemAgent(handler func(Result)) Agent {
	cronScheduler := cron.New()
	cronScheduler.Start()
	agent := &MemAgent{
		Configs:         make(map[util.ID]db.Config),
		schedules:       make(map[util.ID]cron.EntryID),
		cronScheduler:   cronScheduler,
		results:         make(chan Result),
		scheduleStrings: make(map[cron.EntryID]string),
		resultHandler:   handler,
	}
	go agent.sendResults()
	return agent
}

func (a *MemAgent) GetType() AgentType {
	return AT_LOCAL
}

func (a *MemAgent) StoreConfig(c db.Config) error {
	a.Configs[c.Id] = c
	return nil
}

func (a *MemAgent) DeleteConfig(id util.ID) error {
	delete(a.Configs, id)
	return nil
}

func (a *MemAgent) GetAllConfigs() ([]db.Config, error) {
	var configs []db.Config

	return configs, nil
}

func (a *MemAgent) ScheduleConfig(id util.ID, cronString string) error {
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
	a.scheduleStrings[scheduleId] = cronString
	if err != nil {
		log.Error().Msgf("Error scheduling: %v", err)
	} else {
		a.schedules[id] = scheduleId
	}
	return err
}

func (a *MemAgent) UnScheduleConfig(id util.ID) error {
	if _, ok := a.schedules[id]; !ok {
		return errors.New(fmt.Sprintf("No schedule with ID %s", id))
	}
	a.cronScheduler.Remove(a.schedules[id])
	delete(a.scheduleStrings, a.schedules[id])
	return nil
}

func (a *MemAgent) GetAllSchedules() ([]db.Schedule, error) {
	var schedules []db.Schedule

	for k, v := range a.schedules {
		s := db.Schedule{
			Id:       k,
			ConfigId: k,
			Schedule: a.scheduleStrings[v],
		}
		schedules = append(schedules, s)
	}

	return schedules, nil
}

func (a *MemAgent) HandleResult(f func(Result)) {
	a.resultHandler = f
}

func (a *MemAgent) sendResults() {
	for result := range a.results {
		a.resultHandler(result)
	}
}
