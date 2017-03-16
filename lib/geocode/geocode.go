package geocode

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/ryankurte/go-mapbox/lib/base"
	"strings"
)

const (
	apiName    = "geocoding"
	apiVersion = "v5"
	apiMode    = "mapbox.places"
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
	Country      string
	Proximity    base.Location
	Types        []Type
	Autocomplete bool
	BBox         base.BoundingBox
	Limit        uint
}

// ForwardResponse is the response from a forward geocode lookup
type ForwardResponse struct {
	*base.FeatureCollection
	Query []string
}

// Forward geocode lookup
// Finds locations from a place name
func (g *Geocode) Forward(place string, req *ForwardRequestOpts) (*ForwardResponse, error) {

	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	resp := ForwardResponse{}

	err = g.base.Query(apiName, apiVersion, apiMode, strings.Replace(place, " ", "+", -1), &v, &resp)

	return &resp, err
}

// ReverseRequestOpts request options fo reverse geocoding
type ReverseRequestOpts struct {
	Types []Type
	Limit uint
}

// ReverseResponse is the response to a reverse geocode request
type ReverseResponse struct {
	*base.FeatureCollection
	Query []float64
}

// Reverse geocode lookup
// Finds place names from a location
func (g *Geocode) Reverse(loc *base.Location, req *ReverseRequestOpts) (*ReverseResponse, error) {

	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}

	resp := ReverseResponse{}

	query := fmt.Sprintf("%f,%f", loc.Longitude, loc.Latitude)

	err = g.base.Query(apiName, apiVersion, apiMode, query, &v, &resp)

	return &resp, err
}
