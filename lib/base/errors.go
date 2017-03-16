package base

import (
	"errors"
)

// ErrorAPIUnauthorized indicates authorization failed
var ErrorAPIUnauthorized = errors.New("Mapbox API error unauthorized")

// ErrorAPILimitExceeded indicates the API limit has been exceeded
var ErrorAPILimitExceeded = errors.New("Mapbox API error api rate limit exceeded")
