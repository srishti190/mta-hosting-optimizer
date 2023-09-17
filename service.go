package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"mta-hosting-optimizer/dto"
	"net/http"
	"os"
	"strconv"
)

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

	respondWithJSON(w, http.StatusOK, dto.HostnameResponse{
		ResultSet:   result,
		Success:     "True",
		ErrorReason: "",
	})
}

func getInactiveHostNamesForThreshold(threshold int) []string {
	inactiveHosts := make([]string, 0)

	for hostName, ipStatusList := range dto.HostnameMap {
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
	json.NewEncoder(w).Encode(dto.HostnameResponse{
		ResultSet:   nil,
		Success:     "False",
		ErrorReason: message,
	})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
