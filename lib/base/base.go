package base

import (
	"encoding/json"
	"fmt"
	"net/http"
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
}

// NewBase Create a new API base instance
func NewBase(token string) *Base {
	m := &Base{}

	m.token = token

	return m
}

// Query the mapbox API
func (m *Base) Query(api, version, mode, query string, v *url.Values, inst interface{}) error {

	// Add token to args
	v.Set("access_token", m.token)

	// Generate URL
	url := fmt.Sprintf("%s/%s/%s/%s/%s.json", BaseURL, api, version, mode, query)
	fmt.Printf("URL: %s\n", url)

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
