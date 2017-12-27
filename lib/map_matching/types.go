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
)

// MatchingResponse is the response from GetMatching
// https://www.mapbox.com/api-documentation/#match-response-object
type MatchingResponse struct {
	Code       string
	Matchings  []Matchings
	Tracepoint []TracePoint
}

type Coordinate []float64

type GeojsonGeometry struct {
	Coordinates []Coordinate
}

type PolylineGeometry string

// Matchings it a route object with additional confidence field
// https://www.mapbox.com/api-documentation/#match-object
type Matchings struct {
	Confidence float64
	Distance   float64
	Duration   float64
	Geometry   interface{} // Issue: must support polyline (string) or geojson (object)
	Legs       []MatchingLeg
}

func (m *Matchings) GetGeometryGeojson() (*GeojsonGeometry, error) {
	geojson, ok := m.Geometry.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Malformed geojson geometry (expected map[string]interface, received %t)", m.Geometry)
	}

	t, ok := geojson["type"]
	if !ok {
		return nil, fmt.Errorf("Malformed geojson geometry (no type defined)")
	}
	if t != "LineString" {
		return nil, fmt.Errorf("Malformed geojson geometry (incorrect type name: %s)", t)
	}

	v, ok := geojson["coordinates"]
	if !ok {
		return nil, fmt.Errorf("Malformed geojson geometry (no coordinates defined)")
	}
	values, ok := v.([]interface{})
	//values, ok := v.([][]float64)
	if !ok {
		return nil, fmt.Errorf("Malformed geojson geometry (coordinates are not an array of float pairs)")
	}

	geometry := GeojsonGeometry{}
	for _, v := range values {
		value, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("Could not cast value to coordinate slice (type: %t)", v)
		}
		lat, ok := value[0].(float64)
		if !ok {
			return nil, fmt.Errorf("Error casting lat (type: %t)", value[0])
		}
		lng, ok := value[1].(float64)
		if !ok {
			return nil, fmt.Errorf("Error casting lng (type: %t)", value[1])
		}

		geometry.Coordinates = append(geometry.Coordinates, []float64{lat, lng})
	}

	return &geometry, nil
}

func (m *Matchings) GetGeometryPolyline() (string, error) {
	g, ok := m.Geometry.(string)
	if !ok {
		return "", fmt.Errorf("Non polyline geometry (type: %t)", m.Geometry)
	}
	return g, nil
}

//MatchingLeg legs inside the matching object
type MatchingLeg struct {
	Step     []float64
	Summary  string
	Duration float64
	Distance float64
}

// TracePoint represents the location an input point was matched with
type TracePoint struct {
	WaypointIndex  int16
	Location       []float64
	Name           string
	MatchingsIndex int16
}

// OverviewType Type of returned overview geometry
type OverviewType string

const (
	//OverviewFull returns a detailed overview geometry
	OverviewFull OverviewType = "full"
	//OverviewSimplified returns a simplified overview geometry
	OverviewSimplified OverviewType = "simplified"
	//OverviewFalse returns no overview geometry
	OverviewFalse OverviewType = "false"
)

// GeometryType Format of the returned geometry
type GeometryType string

const (
	// GeometryGeojson returns a geojson like geometry
	GeometryGeojson GeometryType = "geojson"
	// GeometryPolyline returns a polyline 5 encoded string like geometry
	GeometryPolyline GeometryType = "polyline"
	// GeometryPolyline6 returns a polyline 6 encode string like geometry
	GeometryPolyline6 GeometryType = "polyline6"
)

// AnnotationType type of metadata to be returned additionally along the route
type AnnotationType string

const (
	// AnnotationDuration returns a additional duration metadata
	AnnotationDuration AnnotationType = "duration"
	// AnnotationDistance returns a additional distance metadata
	AnnotationDistance AnnotationType = "distance"
	// AnnotationSpeed returns a additional speed metadata
	AnnotationSpeed AnnotationType = "speed"
)

// Codes are direction response Codes
// https://www.mapbox.com/api-documentation/#matching-errors
type Codes string

const (
	// CodeOK success response
	CodeOK Codes = "Ok"
	//NoMatchFound No matching route found
	NoMatchFound Codes = "NoMatch"
	//TooManyCoordinates limited to 100 coordinates per request
	TooManyCoordinates Codes = "TooManyCoordinates"
	//CodeProfileNotFound invalid routing profile
	CodeProfileNotFound Codes = "ProfileNotFound"
	// CodeInvalidInput invalid input data to the server
	CodeInvalidInput Codes = "InvalidInput"
)
