package reportconfig

import (
	"context"

	"github.com/rituraj-junglee/container-critic/models"
)

type Repository interface {
	CreateReportConfig(ctx context.Context, report models.ReportConfig) (err error)
	GetReportConfig(ctx context.Context, reportID string) (report models.ReportConfig, err error)
}
