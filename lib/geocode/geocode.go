/**
 * go-mapbox Geocoding Module
 * Wraps the mapbox geocoding API for server side use
 * See https://www.mapbox.com/api-documentation/#geocoding for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package geocode

import (
	"fmt"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/ryankurte/go-mapbox/lib/base"
)

const (
	apiName          = "geocoding"
	apiVersion       = "v5"
	apiMode          = "mapbox.places"
	apiModePermanent = "mapbox.places-permanent"
)

// Type defines geocode location response types
type Type string

const (
	// Country level
	Country Type = "country"
	// Region level
	Region Type = "region"
	// Postcode level
	Postcode Type = "postcode"
	// District level
	District Type = "district"
	// Place level
	Place Type = "place"
	// Locality level
	Locality Type = "locality"
	// Neighborhood level
	Neighborhood Type = "neighborhood"
	// Address level
	Address Type = "address"
	// POI (Point of Interest) level
	POI Type = "poi"
)

// Geocode api wrapper instance
type Geocode struct {
	base *base.Base
}

// NewGeocode Create a new Geocode API wrapper
func NewGeocode(base *base.Base) *Geocode {
	return &Geocode{base}
}

// ForwardRequestOpts request options fo forward geocoding
type ForwardRequestOpts struct {
	Country      string           `url:"country,omitempty"`
	Proximity    []float64        `url:"proximity,omitempty"`
	Types        []Type           `url:"types,omitempty"`
	Language     string           `url:"language,omitempty"`
	Autocomplete bool             `url:"autocomplete,omitempty"`
	BBox         base.BoundingBox `url:"bbox,omitempty"`
	Limit        uint             `url:"limit,omitempty"`
	FuzzyMatch   bool             `url:"fuzzyMatch,omitempty"`
	Routing      bool             `url:"routing,omitempty"`
	worldview    string           `url:"worldview,omitempty"`
}

// ForwardResponse is the response from a forward geocode lookup
type ForwardResponse struct {
	*base.FeatureCollection
	Query []string
}

// Forward geocode lookup
// Finds locations from a place name
func (g *Geocode) Forward(place string, req *ForwardRequestOpts, permanent ...bool) (*ForwardResponse, error) {

	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	resp := ForwardResponse{}

	queryString := strings.Replace(place, " ", "+", -1)
	if len(permanent) > 0 && permanent[0] {
		err = g.base.Query(apiName, apiVersion, apiModePermanent, fmt.Sprintf("%s.json", queryString), &v, &resp)
	} else {
		err = g.base.Query(apiName, apiVersion, apiMode, fmt.Sprintf("%s.json", queryString), &v, &resp)
	}

	return &resp, err
}

// ReverseRequestOpts request options fo reverse geocoding
type ReverseRequestOpts struct {
	Country     string `url:"country,omitempty"`
	Language    string `url:"language,omitempty"`
	ReverseMode string `url:"reverse_mode,omitempty"`
	Routing     bool   `url:"routing,omitempty"`
	Types       []Type `url:"types,omitempty"`
	Limit       uint   `url:"limit,omitempty"`
	Worldview   string `url:"worldview,omitempty"`
}

// ReverseResponse is the response to a reverse geocode request
type ReverseResponse struct {
	*base.FeatureCollection
	Query []float64
}

// Reverse geocode lookup
// Finds place names from a location
func (g *Geocode) Reverse(loc *base.Location, req *ReverseRequestOpts, permanent ...bool) (*ReverseResponse, error) {

	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	resp := ReverseResponse{}

	queryString := fmt.Sprintf("%f,%f.json", loc.Longitude, loc.Latitude)

	if len(permanent) > 0 && permanent[0] {
		err = g.base.Query(apiName, apiVersion, apiModePermanent, queryString, &v, &resp)
	} else {
		err = g.base.Query(apiName, apiVersion, apiMode, queryString, &v, &resp)
	}

	return &resp, err
}
