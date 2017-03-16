/**
 * go-mapbox Base Module
 * Provides a common base for API modules
 * See https://www.mapbox.com/api-documentation/ for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package base

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	// BaseURL Mapbox API base URL
	BaseURL = "https://api.mapbox.com"

	statusRateLimitExceeded = 429
)

// Base Mapbox API base
type Base struct {
	token string
	debug bool
}

// NewBase Create a new API base instance
func NewBase(token string) *Base {
	b := &Base{}

	b.token = token

	return b
}

func (b *Base) SetDebug(debug bool) {
	b.debug = true
}

// Query the mapbox API
func (b *Base) QueryBase(query string, v *url.Values, inst interface{}) error {

	// Add token to args
	v.Set("access_token", b.token)

	// Generate URL
	url := fmt.Sprintf("%s/%s", BaseURL, query)

	if b.debug {
		fmt.Printf("URL: %s\n", url)
	}

	// Create request object
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	request.URL.RawQuery = v.Encode()

	// Create client instance
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if b.debug {
		data, _ := httputil.DumpRequest(request, true)
		fmt.Printf("Request: %s", string(data))
		data, _ = httputil.DumpResponse(resp, true)
		fmt.Printf("Response: %s", string(data))
	}

	if resp.StatusCode == statusRateLimitExceeded {
		return ErrorAPILimitExceeded
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrorAPIUnauthorized
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&inst)
	if err != nil {
		return err
	}

	return nil
}

// Query the mapbox API
func (b *Base) Query(api, version, mode, query string, v *url.Values, inst interface{}) error {

	// Generate URL
	queryString := fmt.Sprintf("%s/%s/%s/%s", api, version, mode, query)

	return b.QueryBase(queryString, v, inst)
}
