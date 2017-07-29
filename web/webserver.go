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
	Scheduler     *controller.AgentScheduler
	ResultCache   map[string][]*messages.AgentProbeResult
	MaxDataPoints map[string]int
}

func NewWebserver() *Webserver {
	ctrl := controller.NewProbeController()
	go ctrl.Start()
	scheduler := controller.NewAgentScheduler(ctrl)
	return &Webserver{
		Scheduler:     scheduler,
		ResultCache:   make(map[string][]*messages.AgentProbeResult),
		MaxDataPoints: make(map[string]int),
	}
}

func (ws *Webserver) Start() {
	ws.Scheduler.Controller.AddResultListener(ws.handleResult)
	log.Info().Msgf("Webserver: starting")
	http.HandleFunc("/api/agents", ws.listAgents)
	http.HandleFunc("/api/agent.control", ws.controlAgent)
	http.HandleFunc("/api/agent.cancel/", ws.cancel)
	http.HandleFunc("/api/results", ws.listResults)
	http.HandleFunc("/api/results/", ws.showResult)

	http.HandleFunc("/", ws.serveFile)
	http.HandleFunc("/html/", ws.serveFile)
	http.HandleFunc("/js/", ws.serveFile)
	http.HandleFunc("/css/", ws.serveFile)
	http.ListenAndServe(":8080", nil)
}

func (ws *Webserver) serveFile(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	filePath := url[1:]
	data, err := ioutil.ReadFile("static/" + filePath)
	log.Info().Msgf("Serving: %s", "static/"+filePath)
	filePathSplit := strings.Split(filePath, ".")
	fileType := filePathSplit[len(filePathSplit)-1]
	switch fileType {
	case "js":
		w.Header().Set("Content-Type", "text/javascript")
	case "html":
		w.Header().Set("Content-Type", "text/html")
	case "css":
		w.Header().Set("Content-Type", "text/css")
	}

	if err != nil {
		return
	}

	w.Write(data)
}

func (ws *Webserver) handleResult(result *messages.AgentProbeResult) {
	ws.ResultCache[result.ResultId] = append(ws.ResultCache[result.ResultId], result)
}

func (ws *Webserver) controlAgent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var request AgentControlRequest
		json.NewDecoder(r.Body).Decode(&request)

		var resultIds []string
		for _, agent := range request.Agents {
			id := util.GenID()
			ws.Scheduler.ScheduleEveryXSeconds(1, agent, &messages.ControllerLatencyRequest{
				Type:       messages.ControllerLatencyRequest_TCP,
				Host:       request.Host,
				ResultId:   id,
				Parameters: request.Options,
			})
			resultIds = append(resultIds, id)
			ws.MaxDataPoints[id] = request.MaxPoints
		}

		response := &AgentControlResponse{}
		response.StatusCode = 0
		response.StatusMessage = "ok"
		response.Results = resultIds

		w.Header().Set("Content-type", "application/json")
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

func (ws *Webserver) listResults(w http.ResponseWriter, r *http.Request) {
	results := ws.ResultCache
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (ws *Webserver) cancel(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	elements := strings.Split(url, "/")
	id := elements[len(elements)-1]
	ws.Scheduler.Cancel(id)
	response := &StandardResponse{}
	response.StatusCode = 0
	response.StatusMessage = "ok"
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ws *Webserver) showResult(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	elements := strings.Split(url, "/")
	id := elements[len(elements)-1]

	results := ws.ResultCache[id]

	response := &ResultResponse{}

	if len(results) > 0 {
		agent := ws.Scheduler.Controller.Agents[results[0].ProbeId]
		response.AgentId = agent.Id
		response.AgentLabel = agent.Info.Label
		response.AgentLocation = agent.Info.Location
		response.TargetHost = results[0].Host
		response.ResultId = results[0].ResultId
		response.Datapoints = make(map[int64]float64)

		if len(results) >= ws.MaxDataPoints[results[0].ResultId] {
			results = results[len(results)-ws.MaxDataPoints[results[0].ResultId]:]
		}

		for _, result := range results {
			var latency time.Duration
			buf := bytes.NewReader(result.Data)
			binary.Read(buf, binary.LittleEndian, &latency)
			response.Datapoints[result.Timestamp] = latency.Seconds() * 1000
		}
	}

	response.StatusCode = 0
	response.StatusMessage = "ok"

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
