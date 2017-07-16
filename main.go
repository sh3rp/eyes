package main

import (
	"os"
	"time"

	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/web"
)

func main() {
	if os.Args[1] == "server" {
		webserver := web.NewWebserver()
		webserver.Start()
	} else if os.Args[1] == "client" {
		a := agent.NewAgent()
		a.Start(os.Args[2])
		for {
			time.Sleep(1 * time.Second)
		}
	}
}
