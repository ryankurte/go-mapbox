/**
 * go-mapbox Mapbox API Modle
 * Wraps the mapbox APIs for golang server (or application) use
 * See https://www.mapbox.com/api-documentation/for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package mapbox

import (
	"github.com/ryankurte/go-mapbox/lib/base"
	"github.com/ryankurte/go-mapbox/lib/directions"
	"github.com/ryankurte/go-mapbox/lib/directions_matrix"
	"github.com/ryankurte/go-mapbox/lib/geocode"
	"github.com/ryankurte/go-mapbox/lib/map_matching"
	"github.com/ryankurte/go-mapbox/lib/maps"
	"github.com/ryankurte/go-mapbox/lib/surface"
)

// Mapbox API Wrapper structure
type Mapbox struct {
	base *base.Base
	// Maps allows fetching of tiles and tilesets
	Maps *maps.Maps
	// Geocode allows forward (by address) and reverse (by lat/lng) geocoding
	Geocode *geocode.Geocode
	// Directions generates directions between arbitrary points
	Directions *directions.Directions
	// Direction Matrix returns all travel times and ways points between multiple points
	DirectionsMatrix *directionsmatrix.DirectionsMatrix
	// MapMatching snaps inaccurate path tracked to a map to produce a clean path
	MapMatching *mapmatching.MapMatching
	// Surface API wrapper
	Surface *surface.Surface
}

// NewMapbox Create a new mapbox API instance
func NewMapbox(token string) *Mapbox {
	m := &Mapbox{}

	// Create base instance
	m.base = base.NewBase(token)

	// Bind modules
	m.Maps = maps.NewMaps(m.base)
	m.Geocode = geocode.NewGeocode(m.base)
	m.Directions = directions.NewDirections(m.base)
	m.DirectionsMatrix = directionsmatrix.NewDirectionsMatrix(m.base)
	m.MapMatching = mapmatching.NewMapMaptching(m.base)
	m.Surface = surface.NewSurface(m.base)

	return m
}
