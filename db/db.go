package db

import (
	"bytes"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/util"
)

var CONFIG_BUCKET = []byte("config")
var SCHEDULE_BUCKET = []byte("schedule")
var ASSIGNMENT_BUCKET = []byte("assignment")

type EyesDB interface {
	SaveConfig(agent.ActionConfig) error
	GetConfig(util.ID) (agent.ActionConfig, error)
	GetConfigs() ([]agent.ActionConfig, error)
	DeleteConfig(util.ID) error

	SaveSchedule(Schedule) error
	GetSchedule(util.ID) (Schedule, error)
	GetSchedules() ([]Schedule, error)
	DeleteSchedule(util.ID) error
}

type BoltEyesDB struct {
	DB *bolt.DB
}

func NewBoltEyesDB(db *bolt.DB) EyesDB {
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(CONFIG_BUCKET)
		tx.CreateBucketIfNotExists(SCHEDULE_BUCKET)
		tx.CreateBucketIfNotExists(ASSIGNMENT_BUCKET)
		return nil
	})
	return &BoltEyesDB{db}
}

type Schedule struct {
	Id            util.ID
	ActioConfigId util.ID
	Schedule      string
}

type Assignment struct {
	Id         util.ID
	AgentId    util.ID
	ScheduleId util.ID
}

func (db BoltEyesDB) SaveConfig(c agent.ActionConfig) error {
	db.put(CONFIG_BUCKET, []byte(c.Id), toJSON(c))
	return nil
}
func (db BoltEyesDB) GetConfig(id util.ID) (agent.ActionConfig, error) {
	data := db.get(CONFIG_BUCKET, []byte(id))
	return toActionConfig(data), nil
}
func (db BoltEyesDB) GetConfigs() ([]agent.ActionConfig, error) {
	var cfgs []agent.ActionConfig
	rows := db.getAll(CONFIG_BUCKET)

	for _, r := range rows {
		cfgs = append(cfgs, toActionConfig(r))
	}
	return cfgs, nil
}
func (db BoltEyesDB) DeleteConfig(id util.ID) error { return nil }

func (db BoltEyesDB) SaveSchedule(s Schedule) error {
	db.put(SCHEDULE_BUCKET, []byte(s.Id), toJSON(s))
	return nil
}
func (db BoltEyesDB) GetSchedule(id util.ID) (Schedule, error) {
	return toSchedule(db.get(SCHEDULE_BUCKET, []byte(id))), nil
}
func (db BoltEyesDB) GetSchedules() ([]Schedule, error) {
	var schedules []Schedule
	rows := db.getAll(SCHEDULE_BUCKET)

	for _, r := range rows {
		schedules = append(schedules, toSchedule(r))
	}
	return schedules, nil
}
func (db BoltEyesDB) DeleteSchedule(id util.ID) error { return nil }

func (db BoltEyesDB) get(table []byte, key []byte) []byte {
	var data []byte
	err := db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(table))
		data = bucket.Get(key)
		return nil
	})

	if err != nil {
		log.Error().Msgf("Error get: %v", err)
	}
	return data
}

func (db BoltEyesDB) getAll(table []byte) [][]byte {
	var dataRows [][]byte

	db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(table))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			dataRows = append(dataRows, v)
		}
		return nil
	})

	return dataRows
}

func (db BoltEyesDB) put(table []byte, key, value []byte) {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(table))
		return bucket.Put(key, value)
	})

	if err != nil {
		log.Error().Msgf("Error put: %v", err)
	}
}

func toJSON(obj interface{}) []byte {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(obj)

	if err != nil {
		log.Error().Msgf("Error encoding to bytes: %v", err)
	}
	return buf.Bytes()
}

func toActionConfig(data []byte) agent.ActionConfig {
	obj := &agent.ActionConfig{}
	buf := &bytes.Buffer{}
	buf.Write(data)
	err := json.NewDecoder(buf).Decode(obj)

	if err != nil {
		log.Error().Msgf("Error decoding to obj: %v", err)
	}
	return *obj
}

func toSchedule(data []byte) Schedule {
	obj := &Schedule{}
	buf := &bytes.Buffer{}
	buf.Write(data)
	err := json.NewDecoder(buf).Decode(obj)

	if err != nil {
		log.Error().Msgf("Error decoding to obj: %v", err)
	}
	return *obj
}

func toAssignment(data []byte) Assignment {
	obj := &Assignment{}
	buf := &bytes.Buffer{}
	buf.Write(data)
	err := json.NewDecoder(buf).Decode(obj)

	if err != nil {
		log.Error().Msgf("Error decoding to obj: %v", err)
	}
	return *obj
}
