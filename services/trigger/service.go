package trigger

import (
	"context"

	"github.com/rituraj-junglee/container-critic/models"
)

//go:generate mwgen -type=Service -template=http -outDir=service -filePrefix=_
type Service interface {
	// http-path /slack/trigger/trivy
	// http-method POST
	TrivyTrigger(ctx context.Context, triggerReq models.TriggerRequest) (err error)
}
