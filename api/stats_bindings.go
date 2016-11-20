package api

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// InitStatsBindings binds the stats endpoints to the router
func InitStatsBindings(ctx context.Context, router *mux.Router, svc StatsService) {
	router.Handle("/v1/stats/{platform}/{region}/{tag}",
		httptransport.NewServer(ctx, func(ctx context.Context, request interface{}) (interface{}, error) {
			req := request.(*getStatsRequest)
			return svc.GetStats(ctx, req)
		}, decodeGetStatsRequest, encodeBasicResponse),
	).Methods(http.MethodGet)
}

// decodeGetStatsRequest is responsible for decoding a new GetStats request
func decodeGetStatsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	v := mux.Vars(r)
	req := &getStatsRequest{
		platform: v["platform"],
		region:   v["region"],
		tag:      v["tag"],
	}
	return req, nil
}
