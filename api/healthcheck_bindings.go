package api

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// InitHealthcheckBindings binds the healthcheck endpoint to the router
func InitHealthcheckBindings(ctx context.Context, router *mux.Router, svc HealthcheckService) {
	router.Handle("/stats/v1/healthcheck",
		httptransport.NewServer(ctx, func(ctx context.Context, req interface{}) (interface{}, error) {
			return svc.GetHealthcheck(ctx, req)
		}, decodeBasicRequest, encodeBasicResponse),
	).Methods(http.MethodGet)
}
