package httpclient

import "net/http"

// Delete calls Delete using the DefaultClient
func Delete(url string) error {
	return DefaultClient.Delete(url, nil)
}

// Delete performs a DELETE request using the passed path
func (c *Client) Delete(path string, headers map[string]string) error {
	// Execute the request and return the response
	_, err := c.Do(NewRequest(http.MethodDelete, path, headers, nil))
	return err
}
