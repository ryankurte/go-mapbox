/**
 * go-mapbox Map Matching Module
 * Wraps the mapbox Map Matching API for server side use
 * See https://www.mapbox.com/api-documentation/#retrieve-a-match for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package mapmatching

// MatchingResponse is the response from GetMatching
// https://www.mapbox.com/api-documentation/#match-response-object
type MatchingResponse struct {
	Code       string
	Matchings  []Matchings
	Tracepoint []TracePoint
}

// Matchings it a route object with additional confidence field
// https://www.mapbox.com/api-documentation/#match-object
type Matchings struct {
	Confidence float64
	Distance   float32
	Duration   float32
	Geometry   string
	Legs       []MatchingLeg
}

//MatchingLeg legs inside the matching object
type MatchingLeg struct {
	Step     []float32
	Summary  string
	Duration float32
	Distance float32
}

// TracePoint represents the location an input point was matched with
type TracePoint struct {
	WaypointIndex  int16
	Location       []float32
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
