package report

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rituraj-junglee/container-critic/models"
	"github.com/rituraj-junglee/container-critic/repo/reportconfig"
	"github.com/rituraj-junglee/container-critic/repo/reportmeta"
)

type service struct {
	timeAdjustedRiskEnabled bool
	reportmetarepo          reportmeta.Repository
	reportconfigrepo        reportconfig.Repository
}

func NewService(timeAdjustedRiskEnabled bool, reportmetarepo reportmeta.Repository, reportconfigrepo reportconfig.Repository) Service {
	return &service{
		timeAdjustedRiskEnabled: timeAdjustedRiskEnabled,
		reportmetarepo:          reportmetarepo,
		reportconfigrepo:        reportconfigrepo,
	}
}

func (s *service) GenerateReport(ctx context.Context, reportReq models.ReportRequest) (reportRes models.ReportResponse, err error) {
	clusterName := reportReq.ClusterName
	targetConfig := make(map[string]int64)

	report := s.generateReport(ctx, reportReq)
	if s.timeAdjustedRiskEnabled {
		prevReport, err := s.reportconfigrepo.GetReportConfig(ctx, clusterName)
		if err != nil {
			log.Println("Error getting report config: ", err)
		}

		if prevReport.ReportID != "" {
			report = s.adjustRiskCost(ctx, report, prevReport)
		}
	}

	temp := s.generateMarkdown(report)

	now := time.Now().UnixMilli()
	for _, result := range report.ReportResults {
		targetConfig[result.Target] = now
	}

	reportRes = models.ReportResponse{
		Template:     temp,
		TotalCost:    report.TotalCost,
		TargetConfig: targetConfig,
	}

	return
}

func (s *service) generateReport(ctx context.Context, req models.ReportRequest) (res models.Report) {
	res.ClusterName = req.ClusterName
	res.Namespace = req.Namespace
	res.Kind = req.Kind

	for _, finding := range req.Findings {
		for _, result := range finding.Results {
			reportResult := models.ReportResult{
				Target:            result.Target,
				TargetsCostData:   s.calculateTargetCostData(ctx, result),
				Vulnerabilities:   result.Vulnerabilities,
				Misconfigurations: result.Misconfigurations,
				Secrets:           result.Secrets,
			}
			res.ReportResults = append(res.ReportResults, reportResult)
		}
	}

	res.TotalCost = s.calculateTotalCost(res.ReportResults)

	return
}

func (s *service) calculateTargetCostData(ctx context.Context, result models.Result) (cd models.TargetCostData) {
	assetValue, _ := s.reportmetarepo.GetAssetValue(ctx, result.Target)
	exploitLikelihoodValue, _ := s.reportmetarepo.GetExploitLikelihood(ctx, result.Target)
	compliancePenaltyValue, _ := s.reportmetarepo.GetCompliancePenalty(ctx, result.Target)
	downtimeCostValue, _ := s.reportmetarepo.GetDowntimeCost(ctx, result.Target)

	var vulnerabilityCost, misConfigCost, secretsCost float64

	for _, vul := range result.Vulnerabilities {
		vulnerabilityCost += calculateVulnerabilityRisk(
			vul.Severity,
			assetValue,
			exploitLikelihoodValue,
			compliancePenaltyValue,
			downtimeCostValue,
		)
	}

	for _, mis := range result.Misconfigurations {
		misConfigCost += calculateMisconfigurationRisk(
			mis.Severity,
			assetValue,
			exploitLikelihoodValue,
			compliancePenaltyValue,
			downtimeCostValue,
		)
	}

	for _, sec := range result.Secrets {
		secretsCost += calculateSecretRisk(
			sec.Severity,
			assetValue,
			exploitLikelihoodValue,
			compliancePenaltyValue,
			downtimeCostValue,
		)
	}

	cd = models.TargetCostData{
		Target:            result.Target,
		VulnerabilityCost: vulnerabilityCost,
		MisconfigCost:     misConfigCost,
		SecretCost:        secretsCost,
		AssetValue:        assetValue,
		ExploitLikelihood: exploitLikelihoodValue,
		CompliancePenalty: compliancePenaltyValue,
		DowntimeCost:      downtimeCostValue,
	}

	return
}

func (s *service) calculateTotalCost(res []models.ReportResult) (totalCost float64) {
	for _, r := range res {
		totalCost += r.TargetsCostData.VulnerabilityCost + r.TargetsCostData.MisconfigCost + r.TargetsCostData.SecretCost
	}
	return
}

func (s *service) generateMarkdown(report models.Report) string {
	var markdown bytes.Buffer

	// Add metadata
	markdown.WriteString(fmt.Sprintf("# Report for Cluster: %s\n\n", report.ClusterName))
	markdown.WriteString(fmt.Sprintf("## Namespace: %s | Kind: %s | Total Cost: %.2f\n\n", report.Namespace, report.Kind, report.TotalCost))

	// Iterate over results
	for _, result := range report.ReportResults {
		markdown.WriteString(fmt.Sprintf("### Target: %s\n\n", result.Target))

		// Target cost table
		markdown.WriteString("| Target | Misconfiguration Cost | Vulnerability Cost | Secrets Cost | Asset Value | ExploitLikelihood | Compliance Penalty | DowntimeCost | CostGrowthRate |\n")
		markdown.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- | --- |\n")
		markdown.WriteString(fmt.Sprintf("| %s | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f | %.2f | \n\n",
			result.TargetsCostData.Target,
			result.TargetsCostData.MisconfigCost,
			result.TargetsCostData.VulnerabilityCost,
			result.TargetsCostData.SecretCost,
			result.TargetsCostData.AssetValue,
			result.TargetsCostData.ExploitLikelihood,
			result.TargetsCostData.CompliancePenalty,
			result.TargetsCostData.DowntimeCost,
			result.TargetsCostData.CostGrowthRate,
		))

		// Misconfigurations table
		markdown.WriteString("#### Misconfigurations\n\n")
		markdown.WriteString("| Title | ID | Description | Message | Resolution | Severity |\n")
		markdown.WriteString("| --- | --- | --- | --- | --- | --- |\n")
		for _, misconfig := range result.Misconfigurations {
			markdown.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
				misconfig.Title, misconfig.ID, misconfig.Description, misconfig.Message, misconfig.Resolution, misconfig.Severity))
		}
		markdown.WriteString("\n")

		// Vulnerabilities table
		markdown.WriteString("#### Vulnerabilities\n\n")
		markdown.WriteString("| Title | VulnerabilityID | Description | PkgName | InstalledVersion | FixedVersion | Severity |\n")
		markdown.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")
		for _, vuln := range result.Vulnerabilities {
			markdown.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |\n",
				vuln.Title, vuln.VulnerabilityID, vuln.Description, vuln.PkgName, vuln.InstalledVersion, vuln.FixedVersion, vuln.Severity))
		}
		markdown.WriteString("\n")
	}

	return markdown.String()
}

func (s *service) adjustRiskCost(ctx context.Context, report models.Report, prevReport models.ReportConfig) models.Report {

	for i, r := range report.ReportResults {
		growthRate, err := s.reportmetarepo.GetGrowthRatePerDay(ctx, r.Target)
		if err != nil {
			log.Println("Error getting growth rate: ", err)
			continue
		}

		if _, ok := prevReport.TargetConfig[r.Target]; !ok {
			continue
		}
		timeElapsed := calculateTimeElapsed(prevReport.TargetConfig[r.Target])

		vulnCost := calculateTimeAdjustedRisk(r.TargetsCostData.VulnerabilityCost, timeElapsed, growthRate)
		misConfigCost := calculateTimeAdjustedRisk(r.TargetsCostData.MisconfigCost, timeElapsed, growthRate)
		secretCost := calculateTimeAdjustedRisk(r.TargetsCostData.SecretCost, timeElapsed, growthRate)

		report.ReportResults[i].TargetsCostData.VulnerabilityCost = vulnCost
		report.ReportResults[i].TargetsCostData.MisconfigCost = misConfigCost
		report.ReportResults[i].TargetsCostData.SecretCost = secretCost
		report.ReportResults[i].TargetsCostData.CostGrowthRate = growthRate
	}

	report.TotalCost = s.calculateTotalCost(report.ReportResults)
	return report
}
