package sampleData

import (
	"mta-hosting-optimizer/dto"
)

func LoadSampleData() error {
	ips := []string{"127.0.0.1", "127.0.0.2", "127.0.0.3", "127.0.0.4", "127.0.0.5", "127.0.0.6"}
	hostNames := []string{"mta-prod-1", "mta-prod-1", "mta-prod-2", "mta-prod-2", "mta-prod-2", "mta-prod-3"}
	activeStatuses := []bool{true, false, true, true, false, false}

	dto.HostnameMap = make(map[string][]dto.HostnameIPActiveStatus)

	for i, ip := range ips {
		hostName := hostNames[i]
		status := activeStatuses[i]

		dto.HostnameMap[hostName] = append(dto.HostnameMap[hostName], dto.HostnameIPActiveStatus{
			IP:     ip,
			Active: status,
		})
	}

	return nil
}
