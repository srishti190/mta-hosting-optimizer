package dto

type HostnameResponse struct {
	ResultSet   []string `json:"ResultSet"`
	Success     string   `json:"Success"`
	ErrorReason string   `json:"ErrorReason"`
}

type HostnameIPActiveStatus struct {
	IP     string
	Active bool
}

var HostnameMap map[string][]HostnameIPActiveStatus
