/**
 * go-mapbox Geocoding Module
 * Wraps the mapbox geocoding API for server side use
 * See https://www.mapbox.com/api-documentation/#geocoding for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package surface

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/ryankurte/go-mapbox/lib/base"
	"strings"
)

const (
	apiName    = "surface"
	apiVersion = "v4"
	mapId      = "mapbox.mapbox-terrain-v2"
)

// Surface api wrapper instance
type Surface struct {
	base *base.Base
}

// NewSurface Create a new Surface API wrapper
func NewSurface(base *base.Base) *Surface {
	return &Surface{base}
}

// RequestOpts request options for the Surface API lookup
type RequestOpts struct {
	Layer           string `url:"layer"`
	Fields          string `url:"fields"`
	GeoJSON         bool   `url:"geojson,omitempty"`
	Points          string `url:"points,omitempty"`
	EncodedPolyline string `url:"encoded_polyline,omitempty"`
	Zoom            uint   `url:"z,omitempty"`
	Interpolate     bool   `url:"interpolate,omitempty"`
}

// DefaultOpts generates a default option structure
func DefaultOpts() RequestOpts {
	return RequestOpts{
		Layer:  "contour",
		Fields: "ele,index",
	}
}

// SetPoints Attaches a sequence of points to the RequestOpts object
func (o *RequestOpts) SetPoints(points []base.Location) {
	lines := make([]string, len(points))
	for i, p := range points {
		lines[i] = fmt.Sprintf("%f,%f", p.Longitude, p.Latitude)
	}
	o.Points = strings.Join(lines, ";")
}

// Response is the response from a surface lookup
type Response struct {
	Results     []Result `json:"results"`
	Attribution string   `json:"attribution"`
}

// Result is a point from the surface API
type Result struct {
	ID        uint          `json:"id"`
	LatLng    base.Location `json:"latlng"`
	Elevation float64       `json:"ele"`
}

// QueryPoints Query with a list of points
func (g *Surface) QueryPoints(locations []base.Location, opts *RequestOpts) (*Response, error) {
	// Attach points to request
	opts.SetPoints(locations)
	// Run Query
	return g.Query(opts)
}

// Query the surface API with the provided RequestOpts
func (g *Surface) Query(opts *RequestOpts) (*Response, error) {

	v, err := query.Values(opts)
	if err != nil {
		return nil, err
	}

	response := Response{}

	queryString := fmt.Sprintf("%s/%s/%s.json", apiVersion, apiName, mapId)

	err = g.base.QueryBase(queryString, &v, &response)

	return &response, err
}
