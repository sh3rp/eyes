package agent

import (
	"fmt"
	"testing"
	"time"

	"github.com/sh3rp/eyes/db"
	"github.com/sh3rp/eyes/util"
	"github.com/stretchr/testify/assert"
)

func TestDummySchedule(t *testing.T) {
	results := []Result{}
	agent := NewMemAgent(func(r Result) {
		results = append(results, r)
	})
	configParameters := make(map[string]string)
	configParameters["a"] = "avalue"
	config := db.Config{
		Id:         util.ID("pants"),
		Action:     A_TEST,
		Parameters: configParameters,
	}
	agent.StoreConfig(config)
	agent.ScheduleConfig(config.Id, "@every 1s")

	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "avalue", results[0].Tags["a"])
	assert.Equal(t, util.ID("pants"), results[0].ConfigId)
}

func TestSSHSchedule(t *testing.T) {
	agent := NewMemAgent(func(r Result) {
		fmt.Printf("[R]: %s\n", string(r.Data))
	})
	configParameters := make(map[string]string)
	configParameters["host"] = "localhost"
	configParameters["username"] = ""
	configParameters["password"] = ""
	configParameters["command"] = "uptime"
	config := db.Config{
		Id:         util.ID("ssh"),
		Action:     A_SSH,
		Parameters: configParameters,
	}
	agent.StoreConfig(config)
	agent.ScheduleConfig(config.Id, "@every 1s")

	time.Sleep(1 * time.Second)
}
