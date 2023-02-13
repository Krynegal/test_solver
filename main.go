package main

import (
	"devices-test/configs"
	"devices-test/workerpool"
	"flag"
	"log"
)

func main() {
	cfg, err := configs.New()
	if err != nil {
		log.Fatal(err)
	}
	flag.IntVar(&cfg.WorkersNum, "w", cfg.WorkersNum, "Server address")
	flag.Parse()
	workerpool.StartWorkerpool(cfg)
}
