package reportmeta

import "context"

type Repository interface {
	GetAssetValue(ctx context.Context, target string) (assetValue float64, err error)
	GetCompliancePenalty(ctx context.Context, target string) (compliancePenalty float64, err error)
	GetDowntimeCost(ctx context.Context, target string) (downtimeCost float64, err error)
	GetExploitLikelihood(ctx context.Context, target string) (exploitLikelihood float64, err error)
	GetGrowthRatePerDay(ctx context.Context, target string) (growthRate float64, err error)
}
