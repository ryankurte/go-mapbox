/**
 * go-mapbox Directions Module Tests
 * Wraps the mapbox directions API for server side use
 * See https://www.mapbox.com/api-documentation/#retrieve-directions for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package directions

import (
	"os"
	"testing"
)

import (
	"github.com/ryankurte/go-mapbox/lib/base"
)

func TestDirections(t *testing.T) {

	b, err := base.NewBase(os.Getenv("MAPBOX_TOKEN"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	Directions := NewDirections(b)

	t.Run("Can Lookup Directions", func(t *testing.T) {
		var opts RequestOpts

		locs := []base.Location{{37.78, -122.42}, {38.91, -77.03}}

		res, err := Directions.GetDirections(locs, RoutingCycling, &opts)
		if err != nil {
			t.Error(err)
		}

		if Codes(res.Code) != CodeOK {
			t.Errorf("Invalid response code: %s", res.Code)
		}

	})

}
