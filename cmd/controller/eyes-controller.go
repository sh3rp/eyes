package main

import (
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/controller"
)

var V_MAJOR = 0
var V_MINOR = 1
var V_PATCH = 0

func main() {
	flags()
	log.Info().Msgf("Net-Eyes controller v%d.%d.%d", V_MAJOR, V_MINOR, V_PATCH)
	c := controller.NewProbeController()
	go c.Start()
	server := controller.NewGRPCServer(9999, c)
	server.Start()
}

func flags() {

}
