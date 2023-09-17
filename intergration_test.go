package main

import (
	"encoding/json"
	"mta-hosting-optimizer/dto"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

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

	var response dto.HostnameResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Success != "True" {
		t.Errorf("Expected Success to be 'True', got '%s'", response.Success)
	}

	expectedHostnames := []string{"mta-prod-1", "mta-prod-2", "mta-prod-3"}
	for _, expected := range expectedHostnames {
		found := false
		for _, result := range response.ResultSet {
			if result == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected hostname '%s' not found in ResultSet", expected)
		}
	}
}
