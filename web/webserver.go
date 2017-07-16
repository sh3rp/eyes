package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/sh3rp/eyes/controller"
)

type Webserver struct {
	Controller *controller.ProbeController
}

func NewWebserver() *Webserver {
	controller := controller.NewProbeController()
	go controller.Start()
	return &Webserver{
		Controller: controller,
	}
}

func (ws *Webserver) Start() {
	log.Printf("Webserver starting")
	http.HandleFunc("/agents", ws.listAgents)
	http.HandleFunc("/agent.control/", ws.controlAgent)
	http.HandleFunc("/agent.test/", ws.testAgent)
	http.ListenAndServe(":8080", nil)
}

func (ws *Webserver) controlAgent(w http.ResponseWriter, r *http.Request) {

}

func (ws *Webserver) listAgents(w http.ResponseWriter, r *http.Request) {
	agents := ws.Controller.Agents
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(agents)
}

func (ws *Webserver) testAgent(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	elements := strings.Split(url, "/")
	id := elements[len(elements)-1]
	log.Printf("Testing agent: %s", id)
	ws.Controller.TestProbe(id)
}
