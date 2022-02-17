package service

import "github.com/labstack/echo/v4"

// newErr creates and returns a new echo HTTPError with the passed status code
// and optional message. Message expected to be of type string or error
func newErr(code int, message ...interface{}) error {
	if len(message) > 0 {
		switch v := message[0].(type) {
		case error:
			return echo.NewHTTPError(code, v.Error())
		case string:
			return echo.NewHTTPError(code, v)
		}
	}
	return echo.NewHTTPError(code, "An error has occurred")
}
