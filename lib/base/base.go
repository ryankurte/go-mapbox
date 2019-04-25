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
	"errors"
	"fmt"
	"io/ioutil"
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
func NewBase(token string) (*Base, error) {
	if token == "" {
		return nil, errors.New("Mapbox API token not found")
	}

	b := &Base{}

	b.token = token

	return b, nil
}

// SetDebug enables debug output for API calls
func (b *Base) SetDebug(debug bool) {
	b.debug = true
}

type MapboxApiMessage struct {
	Message string
}

// QueryRequest make a get with the provided query string and return the response if successful
func (b *Base) QueryRequest(query string, v *url.Values) (*http.Response, error) {
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
		return nil, err
	}
	request.URL.RawQuery = v.Encode()

	// Create client instance
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if b.debug {
		data, _ := httputil.DumpRequest(request, true)
		fmt.Printf("Request: %s", string(data))
		data, _ = httputil.DumpResponse(resp, false)
		fmt.Printf("Response: %s", string(data))
	}

	if resp.StatusCode == statusRateLimitExceeded {
		return nil, ErrorAPILimitExceeded
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrorAPIUnauthorized
	}

	return resp, nil
}

// QueryBase Query the mapbox API and fill the provided instance with the returned JSON
// TODO: Rename this
func (b *Base) QueryBase(query string, v *url.Values, inst interface{}) error {
	// Make request
	resp, err := b.QueryRequest(query, v)
	if err != nil && (resp == nil || resp.StatusCode != http.StatusBadRequest) {
		return err
	}
	defer resp.Body.Close()

	// Read body into buffer
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Handle bad requests with messages
	if resp.StatusCode == http.StatusBadRequest {
		apiMessage := MapboxApiMessage{}
		messageErr := json.Unmarshal(body, &apiMessage)
		if messageErr == nil {
			return fmt.Errorf("api error: %s", apiMessage.Message)
		}
		return fmt.Errorf("Bad Request (400) - no message")
	}

	// Attempt to decode body into inst type
	err = json.Unmarshal(body, &inst)
	if err != nil {
		return err
	}

	return nil
}

// Query the mapbox API
// TODO: Depreciate this
func (b *Base) Query(api, version, mode, query string, v *url.Values, inst interface{}) error {

	// Generate URL
	queryString := fmt.Sprintf("%s/%s/%s/%s", api, version, mode, query)

	return b.QueryBase(queryString, v, inst)
}
