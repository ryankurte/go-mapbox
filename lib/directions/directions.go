/**
 * go-mapbox Directions Module
 * Wraps the mapbox directions API for server side use
 * See https://www.mapbox.com/api-documentation/#retrieve-directions for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package directions

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/ryankurte/go-mapbox/lib/base"
	"strings"
)

const (
	apiName    = "directions"
	apiVersion = "v5"
)

// RoutingProfile defines routing mode for direction finding
type RoutingProfile string

const (
	// RoutingDrivingTraffic mode for automotive routing takes into account current and historic traffic
	RoutingDrivingTraffic RoutingProfile = "mapbox/driving-traffic"
	// RoutingDriving mode for for automovide routing
	RoutingDriving RoutingProfile = "mapbox/driving"
	// RoutingWalking mode for Pedestrian routing
	RoutingWalking RoutingProfile = "mapbox/walking"
	// RoutingCycling mode for bicycle routing
	RoutingCycling RoutingProfile = "mapbox/cycling"
)

type GeometryType string

const (
	GeometryGeojson   GeometryType = "geojson"
	GeometryPolyline  GeometryType = "polyline"
	GeometryPolyline6 GeometryType = "polyline6"
)

type OverviewType string

const (
	OverviewFull       OverviewType = "full"
	OverviewSimplified OverviewType = "simplified"
	OverviewFalse      OverviewType = "false"
)

type AnnotationType string

const (
	AnnotationDuration AnnotationType = "duration"
	AnnotationDistance AnnotationType = "distance"
	AnnotationSpeed    AnnotationType = "speed"
)

type RadiusType string

const (
	RaduisUnlimited RadiusType = "unlimited"
)

// Directions api wrapper instance
type Directions struct {
	base *base.Base
}

// NewDirections Create a new Directions API wrapper
func NewDirections(base *base.Base) *Directions {
	return &Directions{base}
}

// RequestOpts request options for directions api
type RequestOpts struct {
	Alternatives     bool          `url:"alternatives,omitempty"`
	Geometries       *GeometryType `url:"geometries,omitempty"`
	Overview         *OverviewType `url:"overview,omitempty"`
	Radiuses         string        `url:"radiuses,omitempty"`
	Steps            bool          `url:"steps,omitempty"`
	ContinueStraight bool          `url:"continue_straight,omitempty"`
	Bearings         string        `url:"bearings,omitempty"`
}

// SetRadiuses sets radiuses for the maximum distance any coordinate can move when snapped to  nearby road segment.
// This must have the same number of radiuses as locations in the GetDirections request
func (o *RequestOpts) SetRadiuses(radiuses []float64) {
	lines := make([]string, len(radiuses))
	for i, r := range radiuses {
		lines[i] = fmt.Sprintf("%f", r)
	}
	o.Radiuses = strings.Join(lines, ";")
}

// SetBearings builds the bearings query argument from an array of angles and deviations
// Note that this must be used with SetRadiuses and the length of the associated arrays must be the same
func (o *RequestOpts) SetBearings(angles []float64, deviations []float64) error {
	if len(angles) != len(deviations) {
		return fmt.Errorf("RequestOpts.SetBearings error, angle and deviation arrays must have the same length")
	}

	lines := make([]string, len(angles))
	for i := range angles {
		lines[i] = fmt.Sprintf("%f,%f", angles[i], deviations[i])
	}
	o.Bearings = strings.Join(lines, ";")

	return nil
}

// SetAnnotations builds the annotations query argument from an array of annotation types
func (o *RequestOpts) SetAnnotations(annotations []AnnotationType) {
	lines := make([]string, len(annotations))
	for i, a := range annotations {
		lines[i] = fmt.Sprintf("%s", a)
	}
	o.Radiuses = strings.Join(lines, ",")
}

// GetDirections between a set of locations using the specified routing profile
func (g *Directions) GetDirections(locations []base.Location, profile RoutingProfile, opts *RequestOpts) (*DirectionResponse, error) {

	v, err := query.Values(opts)
	if err != nil {
		return nil, err
	}

	coordinateStrings := make([]string, len(locations))
	for i, l := range locations {
		coordinateStrings[i] = fmt.Sprintf("%f,%f", l.Longitude, l.Latitude)
	}
	queryString := strings.Join(coordinateStrings, ";")

	resp := DirectionResponse{}

	err = g.base.Query(apiName, apiVersion, string(profile), queryString, &v, &resp)

	return &resp, err
}
