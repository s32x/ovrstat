package httpclient

import (
	"io/ioutil"
	"net/http"
)

// Response is a basic HTTP response struct containing just the important data
// returned from an HTTP request
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// NewResponse creates a new HTTP Response
func NewResponse(res *http.Response) (*Response, error) {
	// Create a map of all headers
	headers := make(map[string]string)
	for k, v := range res.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	// Decode the body into a slice of bytes
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Return the fully constructed Response
	return &Response{
		StatusCode: res.StatusCode,
		Headers:    headers,
		Body:       body,
	}, nil
}
