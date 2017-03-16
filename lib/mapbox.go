package mapbox

import (
	"github.com/ryankurte/go-mapbox/lib/base"
)

type Mapbox struct {
	base *base.Base
}

func NewMapbox(token string) *Mapbox {
	m := &Mapbox{}

	m.base = base.NewBase(token)

	return m
}
