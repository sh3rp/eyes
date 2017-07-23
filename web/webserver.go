package web

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sh3rp/eyes/controller"
	"github.com/sh3rp/eyes/messages"
	"github.com/sh3rp/eyes/util"
)

type Webserver struct {
	Scheduler   *controller.AgentScheduler
	ResultCache map[string][]*messages.ProbeResult
}

func NewWebserver() *Webserver {
	ctrl := controller.NewProbeController()
	go ctrl.Start()
	scheduler := controller.NewAgentScheduler(ctrl)
	return &Webserver{
		Scheduler:   scheduler,
		ResultCache: make(map[string][]*messages.ProbeResult),
	}
}

func (ws *Webserver) Start() {
	ws.Scheduler.Controller.AddResultListener(ws.handleResult)
	log.Info().Msgf("Webserver: starting")
	http.HandleFunc("/api/agents", ws.listAgents)
	http.HandleFunc("/api/agent.control", ws.controlAgent)
	http.HandleFunc("/api/agent.test/", ws.testAgent)
	http.HandleFunc("/api/results", ws.listResults)
	http.HandleFunc("/api/results/", ws.showResult)

	http.HandleFunc("/", ws.serveFile)
	http.HandleFunc("/js/", ws.serveFile)
	http.HandleFunc("/css/", ws.serveFile)
	http.ListenAndServe(":8080", nil)
}

func (ws *Webserver) serveFile(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	filePath := url[1:]
	data, err := ioutil.ReadFile("static/" + filePath)

	filePathSplit := strings.Split(filePath, ".")
	fileType := filePathSplit[len(filePathSplit)-1]
	log.Info().Msgf("Loading: %s", filePath)
	switch fileType {
	case "js":
		w.Header().Set("Content-Type", "text/javascript")
	case "html":
		w.Header().Set("Content-Type", "text/html")
	case "css":
		w.Header().Set("Content-Type", "text/html")
	}

	if err != nil {
		return
	}

	w.Write(data)
}

func (ws *Webserver) handleResult(result *messages.ProbeResult) {
	ws.ResultCache[result.CmdId] = append(ws.ResultCache[result.CmdId], result)
}

func (ws *Webserver) controlAgent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var request AgentControlRequest
		json.NewDecoder(r.Body).Decode(&request)

		var resultIds []string
		log.Info().Msgf("Agents: %v", request.Agents)
		for _, agent := range request.Agents {
			id := util.GenID()
			ws.Scheduler.ScheduleEveryXSeconds(3, agent, &messages.ProbeCommand{
				Type: messages.ProbeCommand_TCP,
				Host: request.Host,
				Id:   id,
			})
			resultIds = append(resultIds, id)
		}

		response := &AgentControlResponse{}
		response.StatusCode = 0
		response.StatusMessage = "ok"
		response.Results = resultIds

		json.NewEncoder(w).Encode(response)
	} else {
		log.Error().Msgf("POST required for this endpoint")
	}
}

func (ws *Webserver) listAgents(w http.ResponseWriter, r *http.Request) {
	agents := ws.Scheduler.Controller.Agents
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(agents)
}

func (ws *Webserver) testAgent(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	elements := strings.Split(url, "/")
	id := elements[len(elements)-1]
	log.Info().Msgf("Testing agent: %s", id)
	ws.Scheduler.Controller.TestProbe(id)
}

func (ws *Webserver) listResults(w http.ResponseWriter, r *http.Request) {
	results := ws.ResultCache
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (ws *Webserver) showResult(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	elements := strings.Split(url, "/")
	id := elements[len(elements)-1]

	results := ws.ResultCache[id]
	agent := ws.Scheduler.Controller.Agents[results[0].ProbeId]

	response := &ResultResponse{}
	response.AgentId = agent.Id
	response.AgentLabel = agent.Label
	response.AgentLocation = agent.Location
	response.ResultId = results[0].CmdId
	response.Datapoints = make(map[int64]int64)

	for _, result := range results {
		var latency time.Duration
		buf := bytes.NewReader(result.Data)
		binary.Read(buf, binary.LittleEndian, &latency)
		response.Datapoints[result.Timestamp] = int64(latency.Seconds() * 1000)
	}

	response.StatusCode = 0
	response.StatusMessage = "ok"

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
