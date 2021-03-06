package main

import (
	"io/ioutil"
	"errors"
	"github.com/BurntSushi/toml"
	"log"
)

const DEFAULT_CONFIG_LOCATION = "./config.toml"
const CONFIG_FILE_EMPTY = "Configuration file is empty"

type Config struct {
	ReportInterval int

	ServerHost string
	ServerPort uint16

	Logging bool
}

func (c Config) Print() {
	log.Printf("Report Interval: %vms", c.ReportInterval)
	log.Printf("Server address: %v:%v", c.ServerHost, c.ServerPort)
	log.Printf("Logging: %t", c.Logging)
}

func ReadConfig(fileLocation string) (*Config, error) {
	fileDataRaw, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return nil, err
	}
	if len(fileDataRaw) == 0 {
		return nil, errors.New(CONFIG_FILE_EMPTY)
	}
	fileData := string(fileDataRaw)
	var conf Config
	if _, err := toml.Decode(fileData, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func ReadDefaultConfig() (*Config, error) {
	return ReadConfig(DEFAULT_CONFIG_LOCATION)
}
