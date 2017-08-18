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
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type BoundingBox []float64

type Geometry struct {
	Type        string `json:"type"`
	Coordinates Point  `json:"coordinates"`
}

type Context struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	ShortCode string `json:"short_code"`
	WikiData  string `json:"wikidata"`
}

type Properties struct {
	Category string `json:"category"`
	Tel      string `json:"tel"`
	Wikidata string `json:"wikidata"`
	Landmark bool   `json:"landmark"`
	Maki     string `json:"short_code"`
}

type Feature struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Text       string      `json:"text"`
	PlaceName  string      `json:"place_name"`
	PlaceType  []string    `json:"place_type"`
	Relevance  float64     `json:"relevance"`
	Properties Properties  `json:"properties"`
	BBox       BoundingBox `json:"bbox"`
	Center     Point       `json:"center"`
	Geometry   Geometry    `json:"geometry"`
	Context    []Context   `json:"context"`
}

type FeatureCollection struct {
	Type        string    `json:"type"`
	Features    []Feature `json:"features"`
	Attribution string    `json:"attribution"`
}
