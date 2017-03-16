/**
 * go-mapbox Directions Module Types
 * Wraps the mapbox directions API for server side use
 * See https://www.mapbox.com/api-documentation/#retrieve-directions for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package directions

// DirectionResponse is the response from GetDirections
// https://www.mapbox.com/api-documentation/#directions-response-object
type DirectionResponse struct {
	Code      string
	Waypoints []Waypoint
	Routes    []Route
}

// Route A route through (potentially multiple) waypoints.
// https://www.mapbox.com/api-documentation/#route-object
type Route struct {
	Distance float64
	Duration float64
	Geometry string
	Legs     []RouteLeg
}

// Waypoint is an input point snapped to the road network
// https://www.mapbox.com/api-documentation/#waypoint-object
type Waypoint struct {
	Name     string
	Location []float64
}

// RouteLeg A route between two Waypoints
// https://www.mapbox.com/api-documentation/#routeleg-object
type RouteLeg struct {
	Distance   float64
	Duration   float64
	Steps      []RouteStep
	Summary    string
	Annotation Annotation
}

// Annotation conains additional details about each line segment
// https://www.mapbox.com/api-documentation/#routeleg-object
type Annotation struct {
	Distance []float64
	Duration []float64
	Speed    []float64
}

// RouteStep Includes one StepManeuver object and travel to the following RouteStep.
// https://www.mapbox.com/api-documentation/#routestep-object
type RouteStep struct {
	Distance      float64
	Duration      float64
	Geometry      string
	Name          string
	Ref           string
	Destinations  string
	Mode          TransportationMode
	Maneuver      StepManeuver
	Intersections []Intersection
}

// TransportationMode indicates the mode of transportation
// https://www.mapbox.com/api-documentation/#routestep-object
type TransportationMode string

const (
	ModeDriving      TransportationMode = "driving"
	ModeWalking      TransportationMode = "walking"
	ModeFerry        TransportationMode = "ferry"
	ModeCycling      TransportationMode = "cyling"
	ModeUnaccessible TransportationMode = "unaccessible"
)

// Intersection
// https://www.mapbox.com/api-documentation/#routestep-object
type Intersection struct {
	Location []float64
	Bearings []float64
	Entry    []string
	In       uint
	Out      uint
	Lanes    []Lane
}

// Lane
//https://www.mapbox.com/api-documentation/#lane-object
type Lane struct {
	Valid      bool
	Indicatons []string
}

// StepManeuver
// https://www.mapbox.com/api-documentation/#stepmaneuver-object
type StepManeuver struct {
	Locaton       []float64
	BearingBefore float64
	BearingAfter  float64
	Instruction   string
	Type          string
	Modifier      StepModifier
}

// StepModifier indicates the direction change of the maneuver
// https://www.mapbox.com/api-documentation/#stepmaneuver-object
type StepModifier string

const (
	StepModifierUTurn       StepModifier = "uturn"
	StepModifierSharpRight  StepModifier = "sharp right"
	StepModifierRight       StepModifier = "right"
	StepModifierSlightRight StepModifier = "slight right"
	StepModifierStraight    StepModifier = "straight"
	StepModifierSharpLeft   StepModifier = "sharp left"
	StepModifierLeft        StepModifier = "left"
	StepModifierSlightLeft  StepModifier = "slight left"
)

// Codes are direction response Codes
// https://www.mapbox.com/api-documentation/#directions-errors
type Codes string

const (
	CodeOK              Codes = "Ok"
	CodeNoRoute         Codes = "NoRoute"
	CodeNoSegment       Codes = "NoSegment"
	CodeProfileNotFound Codes = "ProfileNotFound"
	CodeInvalidInput    Codes = "InvalidInput"
)
