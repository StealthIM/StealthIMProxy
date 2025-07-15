package main

import (
	"StealthIMProxy/config"
	"StealthIMProxy/conns"
	"StealthIMProxy/service"
	"log"
)

func main() {
	cfg := config.ReadConf()
	log.Printf("Start server [%v]\n", config.Version)

	conns.Init(cfg)

	service.Start(cfg)
}
