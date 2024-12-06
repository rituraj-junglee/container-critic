package models

type ReportRequest struct {
	ClusterName string    `json:"ClusterName"`
	Findings    []Finding `json:"Findings"`
	Date        string    `json:"Date"`
}

type ReportResponse struct {
	Template  string `json:"Template"`
	TotalCost float64
	Date      string
}

type Report struct {
	ReportID      string         `json:"ReportID"`
	ClusterName   string         `json:"ClusterName"`
	Namespace     string         `json:"Namespace"`
	Kind          string         `json:"Kind"`
	ReportResults []ReportResult `json:"ReportResults"`
	TotalCost     float64        `json:"TotalCost"`
	Date          string         `json:"Date"`
}

type ReportResult struct {
	Target            string             `json:"Target"`
	TargetsCostData   TargetCostData     `json:"TargetCostData"`
	Misconfigurations []Misconfiguration `json:"Misconfigurations"`
	Vulnerabilities   []Vulnerability    `json:"Vulnerabilities"`
	Secrets           []Secrets          `json:"Secrets"`
}

type TargetCostData struct {
	Target            string  `json:"Target"`
	VulnerabilityCost float64 `json:"VulnerabilityCost"`
	MisconfigCost     float64 `json:"MisconfigCost"`
	SecretCost        float64 `json:"SecretCost"`
	AssetValue        float64 `json:"AssetValue"`
	ExploitLikelihood float64 `json:"ExploitLikelihood"`
	CompliancePenalty float64 `json:"CompliancePenalty"`
	DowntimeCost      float64 `json:"DowntimeCost"`
}
