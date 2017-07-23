package main

import (
	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/web"
)

var V_MAJOR = 0
var V_MINOR = 1
var V_PATCH = 0

func main() {
	flags()
	log.Info().Msgf("Net-Eyes controller v%d.%d.%d", V_MAJOR, V_MINOR, V_PATCH)
	webserver := web.NewWebserver()
	log.Info().Msgf("Webserver: starting")
	webserver.Start()
}

func flags() {

}
