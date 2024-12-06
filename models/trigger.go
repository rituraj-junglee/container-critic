package models

type TriggerRequest struct {
	ClusterName string    `json:"ClusterName"`
	Namespace   string    `json:"Namespace"`
	Kind        string    `json:"Kind"`
	Findings    []Finding `json:"Findings"`
}

type Finding struct {
	Name      string   `json:"Name"`
	Namespace string   `json:"Namespace"`
	Kind      string   `json:"Kind"`
	Results   []Result `json:"Results"`
}
type Result struct {
	Target            string             `json:"Target"`
	Vulnerabilities   []Vulnerability    `json:"Vulnerabilities"`
	Misconfigurations []Misconfiguration `json:"Misconfigurations"`
	Secrets           []Secrets          `json:"Secrets"`
}

type Vulnerability struct {
	VulnerabilityID  string   `json:"VulnerabilityID"`
	PkgName          string   `json:"PkgName"`
	InstalledVersion string   `json:"InstalledVersion"`
	FixedVersion     string   `json:"FixedVersion"`
	Title            string   `json:"Title"`
	Description      string   `json:"Description"`
	Severity         string   `json:"Severity"`
	References       []string `json:"References"`
}

type Misconfiguration struct {
	Type        string   `json:"Type"`
	ID          string   `json:"ID"`
	AVDID       string   `json:"AVDID"`
	Title       string   `json:"Title"`
	Description string   `json:"Description"`
	Message     string   `json:"Message"`
	Namespace   string   `json:"Namespace"`
	Query       string   `json:"Query"`
	Resolution  string   `json:"Resolution"`
	Severity    string   `json:"Severity"`
	PrimaryURL  string   `json:"PrimaryURL"`
	References  []string `json:"References"`
	Status      string   `json:"Status"`
}

type Secrets struct {
	Severity string `json:"Severity"`
}
