package trigger

import (
	"net/url"

	httpkit "github.com/go-kit/kit/transport/http"
)

func NewHTTPClient(url string, opts ...httpkit.ClientOption) (EndpointsSet, error) {
	u, err := parseUrl(url)
	if err != nil {
		return EndpointsSet{}, err
	}
	return EndpointsSet{
		TrivyTriggerEndpoint: httpkit.NewClient(
			"POST",
			u,
			encodeTrivyTriggerRequest,
			decodeTrivyTriggerResponse,
			opts...,
		).Endpoint(),
	}, nil
}

func parseUrl(urlSt string) (u *url.URL, err error) {
	u, err = url.Parse(urlSt)
	if err != nil {
		return u, err
	}
	return u, nil
}
