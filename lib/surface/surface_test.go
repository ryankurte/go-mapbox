/**
 * go-mapbox Geocoding Module Tests
 * Wraps the mapbox geocoding API for server side use
 * See https://www.mapbox.com/api-documentation/#geocoding for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package surface

import (
	"os"
	"testing"
)

import (
	"github.com/ryankurte/go-mapbox/lib/base"
)

func TestSurface(t *testing.T) {

	token := os.Getenv("MAPBOX_TOKEN")
	if token == "" {
		t.Error("Mapbox API token not found")
		t.FailNow()
	}

	b := base.NewBase(token)
	b.SetDebug(true)

	surface := NewSurface(b)

	t.Run("Can query surface api", func(t *testing.T) {
		opts := DefaultOpts()

		locs := []base.Location{{-122.42, 37.78}, {-77.03, 38.91}}

		res, err := surface.QueryPoints(locs, &opts)
		if err != nil {
			t.Error(err)
		}

		if len(res.Results) == 0 {
			t.Errorf("No results returned")
		}

	})

}
