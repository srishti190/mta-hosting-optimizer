package main

import (
	"log"
	"mta-hosting-optimizer/config"
	"mta-hosting-optimizer/sampleData"
	"net/http"
	"strconv"
)

func main() {
	err := sampleData.LoadSampleData()
	if err != nil {
		log.Fatal("Failed to load ip config sample data: ", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	addr := ":" + strconv.Itoa(cfg.Port)

	http.HandleFunc("/mta-hosting-optimizer", validateThresholdAndGetHostName)
	log.Fatal(http.ListenAndServe(addr, nil))
}
