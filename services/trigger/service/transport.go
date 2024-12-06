package trigger

import (
	"bytes"
	context "context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path"

	"bitbucket.org/junglee_games/go_common/chttp"
	httpkit "github.com/go-kit/kit/transport/http"
	log "github.com/go-kit/log"
	"github.com/gorilla/mux"
	intpkg "github.com/rituraj-junglee/container-critic/services/trigger"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

func MakeHTTPHandler(s intpkg.Service, logger log.Logger, options []httpkit.ServerOption) http.Handler {
	r := mux.NewRouter()

	r.Methods("POST").Path("/slack/trigger/trivy").Handler(httpkit.NewServer(
		MakeTrivyTriggerEndpoint(s),
		decodeTrivyTriggerRequest,
		encodeTrivyTriggerResponse,
		options...,
	))

	return r
}

// Decode Request

func decodeTrivyTriggerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

	var req trivyTriggerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeTrivyTriggerRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods(POST).Path("/slack/trigger/trivy")
	req.Method = "POST"

	req.URL.Path = path.Join(req.URL.Path, "/slack/trigger/trivy")
	return encodeRequest(ctx, req, request)
}

func decodeTrivyTriggerResponse(ctx context.Context, resp *http.Response) (rtModel interface{}, err error) {
	var res trivyTriggerResponse
	err = chttp.DecodeResponse(ctx, resp, &res)
	if err != nil {
		return
	}
	return &res, err
}

func encodeTrivyTriggerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return chttp.EncodeResponse(ctx, w, response)
}

func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(&buf)
	return nil
}
