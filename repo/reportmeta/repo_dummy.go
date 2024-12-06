package reportmeta

import "context"

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) GetAssetValue(_ context.Context, target string) (assetValue float64, err error) {
	return 10000.0, nil
}

func (r *repository) GetCompliancePenalty(_ context.Context, target string) (compliancePenalty float64, err error) {
	return 5000.0, nil
}

func (r *repository) GetDowntimeCost(_ context.Context, target string) (downtimeCost float64, err error) {
	return 2000.0, nil
}

func (r *repository) GetExploitLikelihood(_ context.Context, target string) (exploitLikelihood float64, err error) {
	return 0.8, nil
}

func (r *repository) GetGrowthRatePerDay(_ context.Context, target string) (growthRate float64, err error) {
	return 0.01, nil
}
