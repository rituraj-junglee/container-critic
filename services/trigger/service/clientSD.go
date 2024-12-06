package trigger

import (
	"io"
	"strings"

	csd "bitbucket.org/junglee_games/go_common/sd"
	clb "bitbucket.org/junglee_games/go_common/sd/lb"
	endpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	httpkit "github.com/go-kit/kit/transport/http"
	log "github.com/go-kit/log"
)

func httpClientFactoryMaker(opts ...httpkit.ClientOption) func(string) (EndpointsSet, error) {
	return func(instance string) (EndpointsSet, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance + "/"
		}
		return NewHTTPClient(instance, opts...)
	}
}

var NewHTTPClientSD = sdClientFactory(httpClientFactoryMaker)

func sdClientFactory(
	maker func(opts ...httpkit.ClientOption) func(string) (EndpointsSet, error),
) func(sd.Instancer, sd.Instancer, log.Logger, ...httpkit.ClientOption) EndpointsSet {
	return func(instancer sd.Instancer, activeInstancer sd.Instancer, logger log.Logger, opts ...httpkit.ClientOption) EndpointsSet {
		var endpoints EndpointsSet

		{
			endpointer := csd.NewEndpointer(instancer, trivyTriggerSDFactory(maker(opts...)), logger)
			balancer := clb.NewRoundRobin(endpointer)
			endpoints.TrivyTriggerEndpoint = clb.NoRetry(balancer)
		}

		return endpoints
	}
}

func trivyTriggerSDFactory(clientMaker func(string) (EndpointsSet, error)) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		c, err := clientMaker(instance)
		return c.TrivyTriggerEndpoint, nil, err
	}
}
