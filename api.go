package main

import (
	"strconv"
	"net/http"
	"bytes"
	"errors"
)

type Api struct {
	Config *Config
	HttpClient *http.Client
}

func NewApi(config *Config) *Api {
	return &Api{
		Config: config,
		HttpClient: &http.Client{},
	}
}

func (a Api) Url() string {
	return a.Config.ServerHost + ":" + strconv.Itoa(int(a.Config.ServerPort))
}

func (a Api) SendData(data string) error {
	jsonData := []byte(data)
	req, err := http.NewRequest("POST", a.Url(), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err := a.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Invalid request")
	}
	return nil
}