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
	"github.com/ryankurte/go-mapbox/lib/surface"
)

// Mapbox API Wrapper structure
type Mapbox struct {
	base *base.Base
	// Geocoding API wrapper
	Geocode *geocode.Geocode
	// Directions API wrapper
	Directions *directions.Directions
	// Surface API wrapper
	Surface *surface.Surface
}

// NewMapbox Create a new mapbox API instance
func NewMapbox(token string) *Mapbox {
	m := &Mapbox{}

	// Create base instance
	m.base = base.NewBase(token)

	// Bind modules
	m.Geocode = geocode.NewGeocode(m.base)
	m.Directions = directions.NewDirections(m.base)
	m.Surface = surface.NewSurface(m.base)

	return m
}
