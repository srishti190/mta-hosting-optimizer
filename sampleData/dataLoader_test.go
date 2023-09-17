package sampleData

import (
	"mta-hosting-optimizer/dto"
	"testing"
)

func TestLoadSampleData(t *testing.T) {
	err := LoadSampleData()
	if err != nil {
		panic(err)
	}
	expectedHostNames := []string{"mta-prod-1", "mta-prod-2", "mta-prod-3"}
	for _, hostName := range expectedHostNames {
		_, found := dto.HostnameMap[hostName]
		if !found {
			t.Errorf("Expected hostname '%s' in HostnameMap, but it was not found", hostName)
		}
	}
}
