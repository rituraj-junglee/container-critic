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
	VulnerabilityID  string   `json:"VulnerabilityID" bson:"vulnerability_id"`
	PkgName          string   `json:"PkgName" bson:"pkg_name"`
	InstalledVersion string   `json:"InstalledVersion" bson:"installed_version"`
	FixedVersion     string   `json:"FixedVersion" bson:"fixed_version"`
	Title            string   `json:"Title" bson:"title"`
	Description      string   `json:"Description" bson:"description"`
	Severity         string   `json:"Severity" bson:"severity"`
	References       []string `json:"References" bson:"references"`
}

type Misconfiguration struct {
	Type        string   `json:"Type" bson:"type"`
	ID          string   `json:"ID" bson:"id"`
	AVDID       string   `json:"AVDID" bson:"avd_id"`
	Title       string   `json:"Title" bson:"title"`
	Description string   `json:"Description" bson:"description"`
	Message     string   `json:"Message" bson:"message"`
	Namespace   string   `json:"Namespace" bson:"namespace"`
	Query       string   `json:"Query" bson:"query"`
	Resolution  string   `json:"Resolution" bson:"resolution"`
	Severity    string   `json:"Severity" bson:"severity"`
	PrimaryURL  string   `json:"PrimaryURL" bson:"primary_url"`
	References  []string `json:"References" bson:"references"`
	Status      string   `json:"Status" bson:"status"`
}

type Secrets struct {
	Severity string `json:"Severity" bson:"severity"`
}
