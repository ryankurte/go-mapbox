/**
 * go-mapbox Map Matching Module
 * Wraps the mapbox Map Matching API for server side use
 * See https://www.mapbox.com/api-documentation/#retrieve-a-match for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package mapmatching

import (
	"fmt"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/ryankurte/go-mapbox/lib/base"
)

const (
	apiName    = "matching"
	apiVersion = "v5"
)

// MapMatching api wrapper instance
type MapMatching struct {
	base *base.Base
}

// NewMapMaptching Create a new Map Matching API wrapper
func NewMapMaptching(base *base.Base) *MapMatching {
	return &MapMatching{base}
}

// RoutingProfile defines routing mode for map matching
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

// RequestOpts request options for map matching api
type RequestOpts struct {
	Geometries  GeometryType    `url:"geometries,omitempty"`
	Radiuses    string          `url:"radiuses,omitempty"`
	Steps       bool            `url:"steps,omitempty"`
	Overview    OverviewType    `url:"overview,omitempty"`
	Timestamps  string          `url:"timestamps,omitempty"`
	Annotations *AnnotationType `url:"annotations,omitempty"`
}

// SetRadiuses sets radiuses for the maximum distance any coordinate can move when snapped to nearby road segment.
// This must have the same number of radiuses as locations in the GetMatching request
func (o *RequestOpts) SetRadiuses(radiuses []int) {
	lines := make([]string, len(radiuses))
	for i, r := range radiuses {
		lines[i] = fmt.Sprintf("%v", r)
	}
	o.Radiuses = strings.Join(lines, ";")
}

// SetAnnotations builds the annotations query argument from an array of annotation types
func (o *RequestOpts) SetAnnotations(annotations []AnnotationType) {
	lines := make([]string, len(annotations))
	for i, a := range annotations {
		lines[i] = fmt.Sprintf("%s", a)
	}
	o.Radiuses = strings.Join(lines, ",")
}

// SetTimestamps builds the Timestamps query argument from an array of timestamps types
// This must have the same number of timestamps as locations in the GetMatching request
func (o *RequestOpts) SetTimestamps(timestamps []int64) {
	lines := make([]string, len(timestamps))
	for i, a := range timestamps {
		lines[i] = fmt.Sprintf("%v", a)
	}
	o.Timestamps = strings.Join(lines, ";")
}

// SetGeometries builds the geometry query argument from the specified geometry type
func (o *RequestOpts) SetGeometries(geometrytype GeometryType) {
	o.Geometries = geometrytype
}

// SetOverview builds the overview query argument from the specified overview type
func (o *RequestOpts) SetOverview(overviewtype OverviewType) {
	o.Overview = overviewtype
}

// SetSteps builds the steps query argument from an array of steps option
func (o *RequestOpts) SetSteps(steps bool) {
	o.Steps = steps
}

// GetMatching for a path using the specified routing profile
func (d *MapMatching) GetMatching(path []base.Location, profile RoutingProfile, opts *RequestOpts) (*MatchingResponse, error) {

	v, err := query.Values(opts)
	if err != nil {
		return nil, err
	}

	coordinateStrings := make([]string, len(path))
	for i, l := range path {
		coordinateStrings[i] = fmt.Sprintf("%f,%f", l.Longitude, l.Latitude)
	}
	queryString := strings.Join(coordinateStrings, ";")

	resp := MatchingResponse{}

	err = d.base.Query(apiName, apiVersion, string(profile), queryString, &v, &resp)

	return &resp, err
}
