package report

func calculateVulnerabilityRisk(severity string, assetValue, exploitLikelihood, compliancePenalty, downtimeCost float64) float64 {
	severityFactor := map[string]float64{
		"LOW":      0.1,
		"MEDIUM":   0.3,
		"HIGH":     0.6,
		"CRITICAL": 1.0,
	}[severity]

	return (assetValue * severityFactor) + (exploitLikelihood * compliancePenalty) + downtimeCost
}

func calculateMisconfigurationRisk(criticality string, assetValue, exploitLikelihood, compliancePenalty, downtimeCost float64) float64 {
	criticalityFactor := map[string]float64{
		"LOW":      0.2,
		"MEDIUM":   0.5,
		"HIGH":     0.8,
		"CRITICAL": 1.0,
	}[criticality]

	return (assetValue * criticalityFactor) + (exploitLikelihood * compliancePenalty) + downtimeCost
}

func calculateSecretRisk(sensitivity string, assetValue, exploitPotential, compliancePenalty, incidentCost float64) float64 {
	sensitivityFactor := map[string]float64{
		"LOW":    0.3,
		"MEDIUM": 0.6,
		"HIGH":   0.9,
	}[sensitivity]

	return (assetValue * sensitivityFactor) + (exploitPotential * incidentCost) + compliancePenalty
}

func calculateTimeAdjustedRisk(_ string, baseRisk, timeElapsed, growthRate float64) float64 {
	return baseRisk * (1 + growthRate*timeElapsed)
}
