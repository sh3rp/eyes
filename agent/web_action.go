package agent

import (
	"io/ioutil"
	"net/http"
)

type WebAction struct{}

func (wa *WebAction) Execute(config ActionConfig) (Result, error) {
	res, err := http.Get(config.Parameters["url"])

	if err != nil {
		return Result{}, err
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return Result{}, err
	}

	return Result{Tags: make(map[string]string), DataCode: DATA_OK, ConfigId: config.Id, Id: NewId(), Data: data, Timestamp: Now()}, nil
}
