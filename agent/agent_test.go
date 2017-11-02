package agent

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDummySchedule(t *testing.T) {
	results := []Result{}
	agent := NewMemAgent()
	agent.ResultFuncs = append(agent.ResultFuncs, func(r Result) {
		results = append(results, r)
	})
	configParameters := make(map[string]string)
	configParameters["a"] = "avalue"
	config := ActionConfig{
		Id:         "pants",
		Action:     A_TEST,
		Parameters: configParameters,
	}
	agent.StoreActionConfig(config)
	agent.ScheduleAction(config.Id, "@every 1s")

	time.Sleep(5 * time.Second)

	assert.Equal(t, 5, len(results))
	assert.Equal(t, "avalue", results[0].Tags["a"])
	assert.Equal(t, "pants", results[0].ConfigId)
}

func TestSSHSchedule(t *testing.T) {
	agent := NewMemAgent()
	agent.ResultFuncs = append(agent.ResultFuncs, func(r Result) {
		fmt.Printf("[R]: %s\n", string(r.Data))
	})
	configParameters := make(map[string]string)
	configParameters["host"] = "localhost"
	configParameters["username"] = ""
	configParameters["password"] = ""
	configParameters["command"] = "uptime"
	config := ActionConfig{
		Id:         "ssh",
		Action:     A_SSH,
		Parameters: configParameters,
	}
	agent.StoreActionConfig(config)
	agent.ScheduleAction(config.Id, "@every 1s")

	time.Sleep(10 * time.Second)
}
