package db

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/util"
	"github.com/stretchr/testify/assert"
)

var DUMMY_ID1 = util.ID("dummy1")
var DUMMY_PARAMETERS1 = map[string]string{
	"key1": "value1",
	"key2": "value2",
}

var DUMMY_ID2 = util.ID("dummy2")
var DUMMY_PARAMETERS2 = map[string]string{
	"key3": "value3",
	"key4": "value4",
}

var DUMMY_CONFIG1 = Config{
	Id:         DUMMY_ID1,
	Action:     agent.A_SSH,
	Parameters: DUMMY_PARAMETERS1,
}

var DUMMY_CONFIG2 = Config{
	Id:         DUMMY_ID2,
	Action:     agent.A_SSH,
	Parameters: DUMMY_PARAMETERS2,
}

var DUMMY_SCHEDULE1 = Schedule{
	Id:       DUMMY_ID1,
	ConfigId: DUMMY_ID1,
	Schedule: "@every 1s",
}

var DUMMY_SCHEDULE2 = Schedule{
	Id:       DUMMY_ID2,
	ConfigId: DUMMY_ID2,
	Schedule: "@every 1s",
}

func TestConfig(t *testing.T) {
	db := newDB()
	eyesDB := NewBoltEyesDB(db)

	err := eyesDB.SaveConfig(Config{
		Id:         DUMMY_ID1,
		Action:     agent.A_SSH,
		Parameters: DUMMY_PARAMETERS1,
	})
	assert.Nil(t, err)
	config, err := eyesDB.GetConfig(DUMMY_ID1)
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, DUMMY_ID1, config.Id)
	assert.Equal(t, agent.A_SSH, config.Action)
	assert.Equal(t, DUMMY_PARAMETERS1, config.Parameters)
	rmDB()
}

func TestMultiConfigs(t *testing.T) {
	db := newDB()
	eyesDB := NewBoltEyesDB(db)

	err := eyesDB.SaveConfig(DUMMY_CONFIG1)
	assert.Nil(t, err)
	err = eyesDB.SaveConfig(DUMMY_CONFIG2)
	assert.Nil(t, err)
	configs, err := eyesDB.GetConfigs()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(configs))
	rmDB()
}

func TestSchedule(t *testing.T) {
	db := newDB()
	eyesDB := NewBoltEyesDB(db)

	err := eyesDB.SaveSchedule(DUMMY_SCHEDULE1)
	assert.Nil(t, err)

	schedule, err := eyesDB.GetSchedule(DUMMY_ID1)
	assert.Nil(t, err)
	assert.NotNil(t, schedule)
	assert.Equal(t, DUMMY_ID1, schedule.Id)
	assert.Equal(t, DUMMY_ID1, schedule.ConfigId)
	assert.Equal(t, "@every 1s", schedule.Schedule)
	rmDB()
}

func TestMultiSchdules(t *testing.T) {
	db := newDB()
	eyesDB := NewBoltEyesDB(db)

	err := eyesDB.SaveSchedule(DUMMY_SCHEDULE1)
	assert.Nil(t, err)
	err = eyesDB.SaveSchedule(DUMMY_SCHEDULE2)
	assert.Nil(t, err)

	schedules, err := eyesDB.GetSchedules()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(schedules))
	rmDB()
}

func rmDB() {
	os.Remove("/tmp/test.db")
}

func newDB() *bolt.DB {
	db, err := bolt.Open("/tmp/test.db", 0644, nil)
	if err != nil {
		log.Error().Msg("Error creating db")
	}
	return db
}
