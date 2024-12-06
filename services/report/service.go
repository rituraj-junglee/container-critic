package report

import (
	"context"

	"github.com/rituraj-junglee/container-critic/models"
)

type Service interface {
	GenerateReport(ctx context.Context, reportReq models.ReportRequest) (reportRes models.ReportResponse, err error)
}
