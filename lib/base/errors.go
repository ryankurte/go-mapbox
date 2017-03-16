/**
 * go-mapbox Base Module Errors
 * Defines common errors returned by API modules
 * See https://www.mapbox.com/api-documentation/ for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package base

import (
	"errors"
)

// ErrorAPIUnauthorized indicates authorization failed
var ErrorAPIUnauthorized = errors.New("Mapbox API error unauthorized")

// ErrorAPILimitExceeded indicates the API limit has been exceeded
var ErrorAPILimitExceeded = errors.New("Mapbox API error api rate limit exceeded")
