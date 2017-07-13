package main

import (
	"os"
	"time"

	"github.com/sh3rp/eyes/client"
	"github.com/sh3rp/eyes/controller"
)

func main() {
	if os.Args[1] == "server" {
		c := controller.NewProbeController()
		c.Start()
	} else if os.Args[1] == "client" {
		c := client.ProbeAgent{}
		c.Start("localhost")
		for {
			time.Sleep(1 * time.Second)
		}
	}
}
