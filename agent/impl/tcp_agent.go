package impl

import (
	"crypto/rand"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/oklog/ulid"
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/agent/action"
	"github.com/sh3rp/eyes/msg"
)

type TCPAgent struct {
	Connection     net.Conn
	ControlChannel chan msg.Packet
	ResultChannel  chan agent.Result
	configs        map[ulid.ULID]*agent.ConfigEntry
	configsLock    *sync.Mutex
}

func (a *TCPAgent) SayHello(hello msg.Hello) error { return nil }

func (a *TCPAgent) ConfigureProbe(config msg.ProbeConfig) error {
	id, err := ulid.Parse(config.Id)

	if err != nil {
		return err
	}

	a.configsLock.Lock()
	a.configs[id] = &agent.ConfigEntry{
		Config:  config,
		LastRun: 0,
		Id:      id,
		Active:  false,
	}
	a.configsLock.Unlock()
	return nil
}

func (a *TCPAgent) ActivateProbe(id ulid.ULID) error {
	a.configsLock.Lock()
	defer a.configsLock.Unlock()
	if _, ok := a.configs[id]; ok {
		a.configs[id].Active = true
		return nil
	} else {
		return errors.New("No such configuration id")
	}
}

func (a *TCPAgent) DeactivateProbe(id ulid.ULID) error {
	a.configsLock.Lock()
	defer a.configsLock.Unlock()
	if _, ok := a.configs[id]; ok {
		a.configs[id].Active = false
		return nil
	} else {
		return errors.New("No such configuration id")
	}
}

func (a *TCPAgent) DeleteProbe(id ulid.ULID) error {
	a.configsLock.Lock()
	defer a.configsLock.Unlock()
	delete(a.configs, id)
	return nil
}

func (a *TCPAgent) RunProbeOnce(id ulid.ULID) error {
	a.configsLock.Lock()
	defer a.configsLock.Unlock()
	if _, ok := a.configs[id]; ok {
		entry := a.configs[id]
		go a.execute(ulid.MustNew(uint64(agent.Now()), rand.Reader), entry.Id, entry.Config.Action, entry.Config.Parameters)
		return nil
	} else {
		return errors.New("No such configuration id")
	}
}

func (a *TCPAgent) Start(controllerHost string) {
	log.Info().Msgf("Starting agent")
	for {
		var c net.Conn

		for c == nil {
			c = a.connect(controllerHost)
			time.Sleep(5 * time.Second)
		}

		a.Connection = c

		go a.WriteLoop()

		for {
			data := make([]byte, 4096)
			len, err := c.Read(data)

			if err != nil {
				log.Error().Msgf("ERROR (read): %v", err)
				break
			}

			packet := msg.Packet{}
			err = proto.Unmarshal(data[:len], &packet)

			if err != nil {
				log.Error().Msgf("ERROR (marshal): %v", err)
			} else {
				go a.Dispatch(packet)
			}
		}
	}
}

func (a *TCPAgent) WriteLoop() {
	for {
		select {
		case result := <-a.ResultChannel:
			data, err := proto.Marshal(&msg.ProbeResult{
				Id:        result.ID.String(),
				Data:      result.Data,
				ConfigId:  result.ConfigID.String(),
				Timestamp: result.Timestamp,
			})
			if err != nil {
				log.Error().Msgf("Error (resultChannel): %v", err)
			} else {
				a.Connection.Write(data)
			}
		case control := <-a.ControlChannel:
			data, err := proto.Marshal(&control)
			if err != nil {
				log.Error().Msgf("Error (controlChannel): %v", err)
			} else {
				a.Connection.Write(data)
			}
		}
	}
}

func (a *TCPAgent) Dispatch(message msg.Packet) {
	switch message.Type {
	case msg.Packet_HELLO:
		a.SayHello(*message.GetHello())
	case msg.Packet_PROBE_CONFIG:
		a.ConfigureProbe(*message.GetConfig())
	case msg.Packet_PROBE_ACTION:
		action := *message.GetAction()
		id, err := ulid.Parse(action.Id)

		if err != nil {
			log.Error().Msgf("Error parsing ID for action: %v", err)
			return
		}

		switch action.Action {
		case msg.ProbeAction_ACTIVATE:
			a.ActivateProbe(id)
		case msg.ProbeAction_DEACTIVATE:
			a.DeactivateProbe(id)
		case msg.ProbeAction_DELETE:
			a.DeleteProbe(id)
		case msg.ProbeAction_RUN_ONCE:
			a.RunProbeOnce(id)
		}
	}
}

func (a *TCPAgent) connect(host string) net.Conn {
	c, err := net.Dial("tcp", host+":12121")

	if err != nil {
		log.Error().Msgf("Error connecting: %v", err)
		return nil
	}

	log.Info().Msgf("Connected: %s", host)

	return c
}

func (a *TCPAgent) scheduler() {
	for {
		a.configsLock.Lock()
		for _, v := range a.configs {
			if v.Active && ((v.LastRun + v.Config.Schedule.EveryMilliseconds) < agent.Now()) {
				go a.execute(ulid.MustNew(uint64(agent.Now()), rand.Reader), v.Id, v.Config.Action, v.Config.Parameters)
				v.LastRun = agent.Now()
			}
		}
		a.configsLock.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func (a *TCPAgent) execute(executeId ulid.ULID, configId ulid.ULID, actionLabel string, parameters map[string]string) {
	var theAction agent.Action
	lowerAction := strings.ToLower(actionLabel)
	switch lowerAction {
	case "ssh":
		theAction = &action.SSHExec{}
	case "tcpping":
		theAction = &action.TCPPing{}
	case "snmp":
		theAction = &action.SNMPPoll{}
	}
	result, err := theAction.Execute(executeId, agent.ActionConfig{configId, parameters})
	if err != nil {
		log.Error().Msgf("Error (agent.execute): %v", err)
	} else {
		a.ResultChannel <- result
	}
}
