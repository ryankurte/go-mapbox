/**
 * go-mapbox Base Module Types
 * Provdes common base types for API modles
 * See https://www.mapbox.com/api-documentation/ for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package base


type Point []float64

type Location struct {
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
}

type BoundingBox []float64

type Geometry struct {
	Type        string
	Coordinates Point
}

type Context struct {
	ID        string
	Text      string
	ShortCode string
	WikiData  string
}

type Feature struct {
	ID         string
	Type       string
	Text       string
	PlaceName  string
	PlaceType  []string
	Relevance  float64
	Properties map[string]string
	BBox       BoundingBox
	Center     Point
	Geometry   Geometry
	Context    []Context
}

type FeatureCollection struct {
	Type        string
	Features    []Feature
	Attribution string
}
