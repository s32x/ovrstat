package api

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

// decodeBasicRequest is responsible for decoding a new basic request
func decodeBasicRequest(context.Context, *http.Request) (interface{}, error) {
	return struct{}{}, nil
}

// encodeBasicResponse is responsible for encoding a new basic response
func encodeBasicResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
