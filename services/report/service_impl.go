package report

import (
	"bytes"
	"context"
	"fmt"

	"github.com/rituraj-junglee/container-critic/models"
	"github.com/rituraj-junglee/container-critic/repo/reportmeta"
)

type service struct {
	timeAdjustedRiskEnabled bool
	reportmetarepo          reportmeta.Repository
}

func NewService(reportmetarepo reportmeta.Repository) Service {
	return &service{
		timeAdjustedRiskEnabled: false,
		reportmetarepo:          reportmetarepo,
	}
}

func (s *service) GenerateReport(ctx context.Context, reportReq models.ReportRequest) (reportRes models.ReportResponse, err error) {
	report := s.generateReport(ctx, reportReq)

	temp := s.generateMarkdown(report)

	reportRes = models.ReportResponse{
		Template: temp,
	}

	return
}

// TODO: Implement the service methods
func (s *service) generateReportTemplate(report models.Report) string {
	return fmt.Sprintf("Report for cluster %s\nTotal Cost: %f\n", report.ClusterName, report.TotalCost)
}

func (s *service) generateReport(ctx context.Context, req models.ReportRequest) (res models.Report) {
	res.ClusterName = req.ClusterName

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

	if s.timeAdjustedRiskEnabled {
		growthRatePerDay, _ := s.reportmetarepo.GetGrowthRatePerDay(ctx, result.Target)
		vulnerabilityCost = calculateTimeAdjustedRisk(
			result.Target,
			vulnerabilityCost,
			1,
			growthRatePerDay,
		)
		misConfigCost = calculateTimeAdjustedRisk(
			result.Target,
			misConfigCost,
			1,
			growthRatePerDay,
		)
		secretsCost = calculateTimeAdjustedRisk(
			result.Target,
			secretsCost,
			1,
			growthRatePerDay,
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
	markdown.WriteString(fmt.Sprintf("## Namespace: %s | Kind: %s\n\n", report.Namespace, report.Kind))

	// Iterate over results
	for _, result := range report.ReportResults {
		markdown.WriteString(fmt.Sprintf("### Target: %s\n\n", result.Target))

		// Target cost table
		markdown.WriteString("| Target | Misconfiguration Cost | Vulnerability Cost | Secrets Cost |\n")
		markdown.WriteString("| --- | --- | --- | --- |\n")
		markdown.WriteString(fmt.Sprintf("| %s | %.2f | %.2f | %.2f |\n\n",
			result.TargetsCostData.Target,
			result.TargetsCostData.MisconfigCost,
			result.TargetsCostData.VulnerabilityCost,
			result.TargetsCostData.SecretCost,
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
