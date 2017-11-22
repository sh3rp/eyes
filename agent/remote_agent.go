package agent

import "github.com/sh3rp/eyes/net"

type RemoteAgent struct {
	connection net.Connection
}

func NewRemoteAgent(c net.Connection) Agent {
	return nil
}
