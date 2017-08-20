package controller

import (
	"errors"
	"strconv"

	nats "github.com/nats-io/go-nats"
)

type NATSQueue struct {
	Conn      *nats.Conn
	natsHosts []*NATSServer
}

type NATSServer struct {
	Host string
	Port int
}

func NewNATSQueue(servers []*NATSServer) *NATSQueue {
	return &NATSQueue{
		natsHosts: servers,
	}
}

func (queue *NATSQueue) Connect() error {
	c, err := nats.Connect(queue.generateNATSUrl())
	if err != nil {
		return err
	}
	queue.Conn = c
	return nil
}

func (queue NATSQueue) Subscribe(channel string, f func(data []byte)) error {
	if queue.Conn == nil {
		return errors.New("No connection opened")
	}
	c := queue.Conn
	_, err := c.Subscribe(channel, func(msg *nats.Msg) {
		f(msg.Data)
	})
	if err != nil {
		return err
	}
	return nil
}

func (queue NATSQueue) Publish(channel string, data []byte) error {
	if queue.Conn == nil {
		return errors.New("No connection opened")
	}
	return queue.Conn.Publish(channel, data)
}

func (queue *NATSQueue) generateNATSUrl() string {
	natsURL := ""

	for idx, server := range queue.natsHosts {
		host := "nats://" + server.Host + ":" + strconv.Itoa(server.Port)
		natsURL += host
		if idx < len(queue.natsHosts) {
			natsURL += " "
		}
	}

	return natsURL
}
