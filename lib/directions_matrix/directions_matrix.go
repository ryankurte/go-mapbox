/**
 * go-mapbox Directions Matrix Module
 * Wraps the mapbox directions matrix API for server side use
 * See https://www.mapbox.com/api-documentation/#retrieve-a-matrix for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package directionsmatrix

import (
	"fmt"
	"strings"

	"github.com/JayBusch/go-mapbox/lib/base"
	"github.com/google/go-querystring/query"
)

const (
	apiName    = "directions-matrix"
	apiVersion = "v1"
)

// DirectionsMatrix api wrapper instance
type DirectionsMatrix struct {
	base *base.Base
}

// NewDirectionsMatrix Create a new Directions Matrix API wrapper
func NewDirectionsMatrix(base *base.Base) *DirectionsMatrix {
	return &DirectionsMatrix{base}
}

// RoutingProfile defines routing mode for direction matrix finding
type RoutingProfile string

const (
	// RoutingDriving mode for for automovide routing
	RoutingDriving RoutingProfile = "mapbox/driving"
	// RoutingWalking mode for Pedestrian routing
	RoutingWalking RoutingProfile = "mapbox/walking"
	// RoutingCycling mode for bicycle routing
	RoutingCycling RoutingProfile = "mapbox/cycling"
)

// DirectionMatrixResponse is the response from GetDirections
// https://www.mapbox.com/api-documentation/#matrix-response-format
type DirectionMatrixResponse struct {
	Code         string
	Durations    [][]float64
	Sources      []Waypoint
	Destinations []Waypoint
}

// Waypoint is an input point snapped to the road network
// https://www.mapbox.com/api-documentation/#waypoint-object
type Waypoint struct {
	Name     string
	Location []float64
}

// Codes are direction response Codes
// https://www.mapbox.com/api-documentation/#matrix-errors
type Codes string

const (
	// CodeOK success response
	CodeOK Codes = "Ok"
	//CodeProfileNotFound invalid routing profile
	CodeProfileNotFound Codes = "ProfileNotFound"
	// CodeInvalidInput invalid input data to the server
	CodeInvalidInput Codes = "InvalidInput"
)

// RequestOpts request options for directions api
type RequestOpts struct {
	Sources      string `url:"sources,omitempty"`
	Destinations string `url:"destinations,omitempty"`
}

// SetSources The points which will act as the starting point.
func (o *RequestOpts) SetSources(sources []string) {
	if sources[0] == "all" {
		o.Sources = "all"
	} else {
		lines := make([]string, len(sources))
		for i, r := range sources {
			lines[i] = fmt.Sprintf("%s", r)
		}
		o.Sources = strings.Join(lines, ";")
	}
}

// SetDestinations The points which will act as the destinations.
func (o *RequestOpts) SetDestinations(destinations []string) {
	if destinations[0] == "all" {
		o.Destinations = "all"
	} else {
		lines := make([]string, len(destinations))
		for i, r := range destinations {
			lines[i] = fmt.Sprintf("%s", r)
		}
		o.Destinations = strings.Join(lines, ";")
	}
}

// GetDirectionsMatrix between a set of locations using the specified routing profile
func (d *DirectionsMatrix) GetDirectionsMatrix(locations []base.Location, profile RoutingProfile, opts *RequestOpts) (*DirectionMatrixResponse, error) {

	v, err := query.Values(opts)
	if err != nil {
		return nil, err
	}

	coordinateStrings := make([]string, len(locations))
	for i, l := range locations {
		coordinateStrings[i] = fmt.Sprintf("%f,%f", l.Longitude, l.Latitude)
	}
	queryString := strings.Join(coordinateStrings, ";")

	resp := DirectionMatrixResponse{}

	err = d.base.Query(apiName, apiVersion, string(profile), queryString, &v, &resp)

	return &resp, err
}
