package models

type ReportRequest struct {
	ClusterName string    `json:"ClusterName"`
	Namespace   string    `json:"Namespace"`
	Kind        string    `json:"Kind"`
	Findings    []Finding `json:"Findings"`
}

type ReportResponse struct {
	Template     string           `json:"Template"`
	TotalCost    float64          `json:"TotalCost"`
	TargetConfig map[string]int64 `json:"TargetConfig"`
}

type Report struct {
	ClusterName   string         `json:"ClusterName" bson:"cluster_name"`
	Namespace     string         `json:"Namespace" bson:"namespace"`
	Kind          string         `json:"Kind" bson:"kind"`
	ReportResults []ReportResult `json:"ReportResults" bson:"report_results"`
	TotalCost     float64        `json:"TotalCost" bson:"total_cost"`
}

type ReportResult struct {
	Target            string             `json:"Target" bson:"target"`
	TargetsCostData   TargetCostData     `json:"TargetCostData" bson:"target_cost_data"`
	Misconfigurations []Misconfiguration `json:"Misconfigurations" bson:"misconfigurations"`
	Vulnerabilities   []Vulnerability    `json:"Vulnerabilities" bson:"vulnerabilities"`
	Secrets           []Secrets          `json:"Secrets" bson:"secrets"`
}

type TargetCostData struct {
	Target            string  `json:"Target" bson:"target"`
	VulnerabilityCost float64 `json:"VulnerabilityCost" bson:"vulnerability_cost"`
	MisconfigCost     float64 `json:"MisconfigCost" bson:"misconfig_cost"`
	SecretCost        float64 `json:"SecretCost" bson:"secret_cost"`
	AssetValue        float64 `json:"AssetValue" bson:"asset_value"`
	ExploitLikelihood float64 `json:"ExploitLikelihood" bson:"exploit_likelihood"`
	CompliancePenalty float64 `json:"CompliancePenalty" bson:"compliance_penalty"`
	DowntimeCost      float64 `json:"DowntimeCost" bson:"downtime_cost"`
	CostGrowthRate    float64 `json:"CostGrowthRate" bson:"cost_growth_rate"`
}

type ReportConfig struct {
	ReportID     string           `json:"ReportID" bson:"_id"`
	TargetConfig map[string]int64 `json:"TargetConfig" bson:"target_config"`
	Timestamp    int64            `json:"Timestamp" bson:"timestamp"`
}
