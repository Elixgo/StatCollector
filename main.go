package main

import (
	"log"
)

func main() {
	log.Println("#==> Starting Elixgo data collector")
	config, err := ReadDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	config.Print()
}
