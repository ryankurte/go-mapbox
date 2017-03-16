package mapbox

import (
	"github.com/ryankurte/go-mapbox/lib/base"
	"github.com/ryankurte/go-mapbox/lib/geocode"
)

// Mapbox API Wrapper structure
type Mapbox struct {
	base    *base.Base
	Geocode *geocode.Geocode
}

// NewMapbox Create a new mapbox API instance
func NewMapbox(token string) *Mapbox {
	m := &Mapbox{}

	// Create base instance
	m.base = base.NewBase(token)

	// Bind modules
	m.Geocode = geocode.NewGeocode(m.base)

	return m
}
