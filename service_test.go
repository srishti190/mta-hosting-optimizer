package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := loadSampleData()
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestLoadSampleData(t *testing.T) {
	expectedHostNames := []string{"mta-prod-1", "mta-prod-2", "mta-prod-3"}
	for _, hostName := range expectedHostNames {
		_, found := HostnameMap[hostName]
		if !found {
			t.Errorf("Expected hostname '%s' in ipMap, but it was not found", hostName)
		}
	}
}

func TestGetInactiveHostNamesForThreshold(t *testing.T) {
	threshold := 1
	expectedResult := []string{"mta-prod-1", "mta-prod-3"}

	result := getInactiveHostNamesForThreshold(threshold)

	if len(result) != len(expectedResult) {
		t.Errorf("Expected %d inefficient instances, but got %d", len(expectedResult), len(result))
	}
}

func TestGetHostNameForValidRequest(t *testing.T) {
	os.Setenv("X", "1")
	defer os.Unsetenv("X")

	req := httptest.NewRequest(http.MethodGet, "/mta-hosting-optimizer", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	validateThresholdAndGetHostName(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var response HostnameResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response body: %s", err)
	}

	expectedResponse := []string{"mta-prod-1", "mta-prod-3"}
	if len(response.Data) != len(expectedResponse) {
		t.Errorf("Expected %d result(s), but got %d", len(expectedResponse), len(response.Data))
	}
}

func TestGetHostNameForInvalidThresholdRequest(t *testing.T) {
	os.Setenv("X", "invalid_threshold")
	defer os.Unsetenv("X")

	req := httptest.NewRequest(http.MethodGet, "/mta-hosting-optimizer", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	validateThresholdAndGetHostName(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, but got %d", http.StatusInternalServerError, w.Code)
	}

	var response HostnameResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response body: %s", err)
	}

	if response.Success != "False" {
		t.Errorf("Expected status 'Error', but got '%s'", response.Success)
	}

	expectedErrorMessage := "Error converting string to int"
	if response.Reason != expectedErrorMessage {
		t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMessage, response.Reason)
	}

	if response.Data != nil {
		t.Errorf("Expected result data as 'Nil', but got '%s'", response.Data)
	}
}

func TestGetHostNameForMissingThresholdRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/mta-hosting-optimizer", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	validateThresholdAndGetHostName(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var response HostnameResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response body: %s", err)
	}

	expectedResponse := []string{"mta-prod-1", "mta-prod-3"}
	if len(response.Data) != len(expectedResponse) {
		t.Errorf("Expected %d result(s), but got %d", len(expectedResponse), len(response.Data))
	}
}

func TestGetEnv(t *testing.T) {
	testCases := []struct {
		Name           string
		Key            string
		DefaultValue   string
		EnvValue       string
		ExpectedResult string
	}{
		{
			Name:           "Environment variable exists",
			Key:            "EXISTING_VARIABLE",
			DefaultValue:   "default_value",
			EnvValue:       "custom_value",
			ExpectedResult: "custom_value",
		},
		{
			Name:           "Environment variable missing",
			Key:            "MISSING_VARIABLE",
			DefaultValue:   "default_value",
			EnvValue:       "",
			ExpectedResult: "default_value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.EnvValue != "" {
				os.Setenv(tc.Key, tc.EnvValue)
				defer os.Unsetenv(tc.Key) // Clean up after the test
			}

			result := getEnv(tc.Key, tc.DefaultValue)

			if result != tc.ExpectedResult {
				t.Errorf("Expected result '%s', but got '%s'", tc.ExpectedResult, result)
			}
		})
	}
}

func TestGoDotEnvVariable1(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	result := GoDotEnvVariable("TEST_ENV_VAR")

	expected := "test_value"
	if result != expected {
		t.Errorf("Expected result '%s', but got '%s'", expected, result)
	}
}

func TestIntegration(t *testing.T) {
	os.Setenv("X", "2")
	defer os.Unsetenv("X")

	req, err := http.NewRequest("GET", "/mta-hosting-optimizer", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(validateThresholdAndGetHostName)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
