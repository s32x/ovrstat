package httpclient

// httpclient is a convenience package for executing HTTP requests. It's safe
// in that it always closes response bodies and returns byte slices, strings or
// decodes responses into interfaces

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

// Client is an http.Client wrapper
type Client struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

// DefaultClient is a default Client for using without having to declare a
// Client
var DefaultClient = NewBaseClient()

// NewBaseClient creates a new Client reference given a client timeout
func NewBaseClient() *Client {
	return &Client{client: &http.Client{}}
}

// SetTimeout sets the timeout on the httpclients client
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.client.Timeout = timeout
	return c
}

// SetBaseURL sets the baseURL on the Client which will be used on all
// subsequent requests
func (c *Client) SetBaseURL(url string) *Client {
	c.baseURL = url
	return c
}

// SetHeaders sets the headers on the Client which will be used on all
// subsequent requests
func (c *Client) SetHeaders(headers map[string]string) *Client {
	c.headers = headers
	return c
}

// Do performs the request and returns a fully populated Response
func (c *Client) Do(req *Request) (*Response, error) {
	// Build the full request URL
	url := c.baseURL + req.Path

	// Encode the body if one was passed
	var b io.ReadWriter
	if req.Body != nil {
		b = bytes.NewBuffer(req.Body)
	}

	// Generate a new request using the new URL
	r, err := http.NewRequest(req.Method, url, b)
	if err != nil {
		return nil, err
	}

	// Add any client and passed headers to the new request
	if c.headers != nil {
		for k, v := range c.headers {
			r.Header.Set(k, v)
		}
	}
	if req.Headers != nil {
		for k, v := range req.Headers {
			r.Header.Set(k, v)
		}
	}

	// Execute the fully constructed request
	res, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Decode the response into a Response and return
	return NewResponse(res)
}
