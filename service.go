package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type HostnameResponse struct {
	Data    []string `json:"data"`
	Success string   `json:"success"`
	Reason  string   `json:"reason"`
}

type HostnameIPActiveStatus struct {
	IP     string
	Active bool
}

var HostnameMap map[string][]HostnameIPActiveStatus

func main() {
	err := loadSampleData()
	if err != nil {
		log.Fatal("Failed to load mock data: ", err)
	}

	http.HandleFunc("/mta-hosting-optimizer", validateThresholdAndGetHostName)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func loadSampleData() error {
	ips := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"}
	hostNames := []string{"mta-prod-1", "mta-prod-1", "mta-prod-2", "mta-prod-2", "mta-prod-2", "mta-prod-3"}
	activeStatuses := []bool{true, false, true, true, false, false}

	HostnameMap = make(map[string][]HostnameIPActiveStatus)

	for i, ip := range ips {
		hostName := hostNames[i]
		status := activeStatuses[i]

		HostnameMap[hostName] = append(HostnameMap[hostName], HostnameIPActiveStatus{
			IP:     ip,
			Active: status,
		})
	}

	return nil
}

func validateThresholdAndGetHostName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	thresholdEnv := GoDotEnvVariable("X")
	threshold, err := strconv.Atoi(thresholdEnv)

	if err != nil {
		errorMsg := "Error converting string to int"
		log.Printf("%s: %v\n", errorMsg, err)
		respondWithError(w, http.StatusInternalServerError, errorMsg)
		return
	}

	result := getInactiveHostNamesForThreshold(threshold)

	respondWithJSON(w, http.StatusOK, HostnameResponse{
		Data:    result,
		Success: "True",
		Reason:  "",
	})
}

func getInactiveHostNamesForThreshold(threshold int) []string {
	inactiveHosts := make([]string, 0)

	for hostName, ipStatusList := range HostnameMap {
		activeIPCount := 0

		for _, ipStatus := range ipStatusList {
			if ipStatus.Active {
				activeIPCount++
			}
		}
		if activeIPCount <= threshold {
			inactiveHosts = append(inactiveHosts, hostName)
		}
	}

	return inactiveHosts
}

func GoDotEnvVariable(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return getEnv(key, "1")
}

func getEnv(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return value
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(HostnameResponse{
		Data:    nil,
		Success: "False",
		Reason:  message,
	})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
