package report

import (
	"context"

	"github.com/rituraj-junglee/container-critic/models"
)

type Repository interface {
	CreateReport(ctx context.Context, report models.Report) (err error)
	GetReport(ctx context.Context, reportID string) (report models.Report, err error)
	GetReports(ctx context.Context) (reports []models.Report, err error)
}
