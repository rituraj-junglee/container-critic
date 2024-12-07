package trigger

import (
	"context"

	"github.com/rituraj-junglee/container-critic/models"
)

// @microgen http

type Service interface {
	// @http-method POST
	// @http-path  slack/trigger/trivy
	TrivyTrigger(ctx context.Context, triggerReq models.TriggerRequest) (err error)
}
