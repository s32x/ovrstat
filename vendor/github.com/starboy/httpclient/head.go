package httpclient

import "net/http"

// Head calls Head using the DefaultClient
func Head(url string) error {
	return DefaultClient.Head(url, nil)
}

// Head performs a HEAD request using the passed path
func (c *Client) Head(path string, headers map[string]string) error {
	// Execute the request and return the response
	_, err := c.Do(NewRequest(http.MethodHead, path, headers, nil))
	return err
}
