package base

import (
	"github.com/ryankurte/go-geojson"
)

type FeatureCollection2 geojson.FeatureCollection

type Point []float64

type Location struct {
	Longitude float64
	Latitude  float64
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
