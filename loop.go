package main

import (
	"log"
	"time"
)

func loop(config *Config) {
	var data, err = NewData(config)
	var api = NewApi(config)
	if err != nil {
		log.Fatal(err)
	}
	data.UpdateData()
	for {
		data.UpdateData()
		json, err := data.ToJsonString()
		if err != nil {
			log.Fatal(err)
		}
		err = api.SendData(json)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Millisecond * time.Duration(config.ReportInterval))
	}
}
