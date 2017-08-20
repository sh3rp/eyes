package config

import (
	"encoding/json"
	"os"
)

type ClientConfig struct {
	ControllerHosts []*ControllerConfig `json:"hosts"`
}

type ControllerConfig struct {
	Label string `json:"label"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
}

func (c *ClientConfig) Read(filename string) error {
	file, err := os.Open(filename)

	if err != nil {
		return err
	}

	return json.NewDecoder(file).Decode(c)
}
