package db

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/util"
)

var CONFIG_BUCKET = []byte("config")
var SCHEDULE_BUCKET = []byte("schedule")
var DEPLOYMENT_BUCKET = []byte("deployment")
var AGENT_BUCKET = []byte("agent")

type EyesDB interface {
	SaveConfig(Config) error
	GetConfig(util.ID) (Config, error)
	GetConfigs() ([]Config, error)
	DeleteConfig(util.ID) error

	SaveSchedule(Schedule) error
	GetSchedule(util.ID) (Schedule, error)
	GetSchedules() ([]Schedule, error)
	DeleteSchedule(util.ID) error

	SaveDeployment(Deployment) error
	GetDeployment(util.ID) (Deployment, error)
	GetDeployments() ([]Deployment, error)
	DeleteDeployment(util.ID) error

	SaveAgent(Agent) error
	GetAgent(util.ID) (Agent, error)
	GetAgents() ([]Agent, error)
	DeleteAgent(util.ID) error
}

type BoltEyesDB struct {
	DB *bolt.DB
}

func NewBoltEyesDB(db *bolt.DB) EyesDB {
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(CONFIG_BUCKET)
		tx.CreateBucketIfNotExists(SCHEDULE_BUCKET)
		tx.CreateBucketIfNotExists(DEPLOYMENT_BUCKET)
		tx.CreateBucketIfNotExists(AGENT_BUCKET)
		return nil
	})
	return &BoltEyesDB{db}
}

type Config struct {
	Id         util.ID
	Action     int
	Parameters map[string]string
}

type Schedule struct {
	Id       util.ID
	ConfigId util.ID
	Schedule string
}

type Deployment struct {
	Id         util.ID
	Agents     []util.ID
	ScheduleId util.ID
}

type Agent struct {
	Id        util.ID
	IpAddress string
}

func (db BoltEyesDB) SaveConfig(c Config) error {
	db.put(CONFIG_BUCKET, []byte(c.Id), toJSON(c))
	return nil
}
func (db BoltEyesDB) GetConfig(id util.ID) (Config, error) {
	data := db.get(CONFIG_BUCKET, []byte(id))
	return toConfig(data), nil
}
func (db BoltEyesDB) GetConfigs() ([]Config, error) {
	var cfgs []Config
	rows := db.getAll(CONFIG_BUCKET)

	for _, r := range rows {
		cfgs = append(cfgs, toConfig(r))
	}
	return cfgs, nil
}
func (db BoltEyesDB) DeleteConfig(id util.ID) error {
	return db.delete(CONFIG_BUCKET, []byte(id))
}

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
func (db BoltEyesDB) DeleteSchedule(id util.ID) error {
	return db.delete(SCHEDULE_BUCKET, []byte(id))
}

func (db BoltEyesDB) SaveDeployment(deployment Deployment) error {
	db.put(DEPLOYMENT_BUCKET, []byte(deployment.Id), toJSON(deployment))
	return nil
}

func (db BoltEyesDB) GetDeployment(id util.ID) (Deployment, error) {
	return toDeployment(db.get(DEPLOYMENT_BUCKET, []byte(id))), nil
}

func (db BoltEyesDB) GetDeployments() ([]Deployment, error) {
	var deployments []Deployment
	rows := db.getAll(DEPLOYMENT_BUCKET)

	for _, r := range rows {
		deployments = append(deployments, toDeployment(r))
	}
	return deployments, nil
}

func (db BoltEyesDB) DeleteDeployment(id util.ID) error {
	return db.delete(DEPLOYMENT_BUCKET, []byte(id))
}

func (db BoltEyesDB) SaveAgent(agent Agent) error {
	db.put(AGENT_BUCKET, []byte(agent.Id), toJSON(agent))
	return nil
}
func (db BoltEyesDB) GetAgent(id util.ID) (Agent, error) {
	return toAgent(db.get(AGENT_BUCKET, []byte(id))), nil
}

func (db BoltEyesDB) GetAgents() ([]Agent, error) {
	var agents []Agent
	rows := db.getAll(AGENT_BUCKET)

	for _, r := range rows {
		agents = append(agents, toAgent(r))
	}
	return agents, nil
}

func (db BoltEyesDB) DeleteAgent(id util.ID) error {
	return db.delete(AGENT_BUCKET, []byte(id))
}

// helpers

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

func (db BoltEyesDB) delete(table []byte, key []byte) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(table)
		b.Delete(key)
		return nil
	})
}

func toJSON(obj interface{}) []byte {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(obj)

	if err != nil {
		log.Error().Msgf("Error encoding to bytes: %v", err)
	}
	return buf.Bytes()
}

func toConfig(data []byte) Config {
	obj := &Config{}
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

func toDeployment(data []byte) Deployment {
	obj := &Deployment{}
	buf := &bytes.Buffer{}
	buf.Write(data)
	err := json.NewDecoder(buf).Decode(obj)

	if err != nil {
		log.Error().Msgf("Error decoding to obj: %v", err)
	}
	return *obj
}

func toAgent(data []byte) Agent {
	obj := &Agent{}
	buf := &bytes.Buffer{}
	buf.Write(data)
	err := json.NewDecoder(buf).Decode(obj)

	if err != nil {
		log.Error().Msgf("Error decoding to obj: %v", err)
	}
	return *obj
}

func NewDB(dir, name string) *bolt.DB {
	db, err := bolt.Open(fmt.Sprintf("%s/%s.db", dir, name), 0644, nil)
	if err != nil {
		log.Error().Msg("Error creating db")
	}
	return db
}
