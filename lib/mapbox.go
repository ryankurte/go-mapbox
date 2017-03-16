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
	"github.com/ryankurte/go-mapbox/lib/geocode"
	"github.com/ryankurte/go-mapbox/lib/maps"
)

// Mapbox API Wrapper structure
type Mapbox struct {
	base       *base.Base
	Maps       *maps.Maps
	Geocode    *geocode.Geocode
	Directions *directions.Directions
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

	return m
}
