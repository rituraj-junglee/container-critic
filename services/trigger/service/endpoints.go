package trigger

import (
	context "context"

	endpoint "github.com/go-kit/kit/endpoint"
	models "github.com/rituraj-junglee/container-critic/models"
	intpkg "github.com/rituraj-junglee/container-critic/services/trigger"
)

type EndpointsSet struct {
	TrivyTriggerEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s intpkg.Service) EndpointsSet {
	return EndpointsSet{
		TrivyTriggerEndpoint: MakeTrivyTriggerEndpoint(s),
	}
}

type trivyTriggerRequest struct {
	TriggerReq models.TriggerRequest `json:"triggerReq"`
}

type trivyTriggerResponse struct {
}

func MakeTrivyTriggerEndpoint(s intpkg.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(trivyTriggerRequest)

		err = s.TrivyTrigger(ctx, req.TriggerReq)

		return trivyTriggerResponse{}, err
	}
}

func (e EndpointsSet) TrivyTrigger(ctx context.Context, triggerReq models.TriggerRequest) (err error) {
	request := trivyTriggerRequest{
		TriggerReq: triggerReq,
	}

	_, err = e.TrivyTriggerEndpoint(ctx, request)

	if err != nil {
		return
	}

	return err
}
