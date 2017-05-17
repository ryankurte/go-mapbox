/**
 * go-mapbox Maps Module Tests
 * Wraps the mapbox Maps API for server side use
 * See https://www.mapbox.com/api-documentation/#maps for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package maps

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryankurte/go-mapbox/lib/base"
)

const delta = 1e-6

func TestMercator(t *testing.T) {
	zoom := uint64(4)
	size := uint64(256)
	fsize := float64(size)

	loc := base.Location{-45.942805, 166.568500}

	t.Run("Performs mercator projections to global pixels", func(t *testing.T) {
		x, y := MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom, size)
		assert.EqualValues(t, 15.0, math.Floor(x/fsize))
		assert.EqualValues(t, 10.0, math.Floor(y/fsize))

		// Increase zoom scale x2 multiplies location by 4
		x, y = MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom+2, size)
		assert.EqualValues(t, 15.0*4+1, math.Floor(x/fsize))
		assert.EqualValues(t, 10.0*4+1, math.Floor(y/fsize))

		// Doubling tile size doubles pixel location
		x, y = MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom, size*2)
		assert.EqualValues(t, 15.0*2, math.Floor(x/fsize))
		assert.EqualValues(t, 10.0*2, math.Floor(y/fsize))
	})

	t.Run("Performs mercator projections to tile IDs", func(t *testing.T) {
		x, y := MercatorLocationToTileID(loc.Latitude, loc.Longitude, zoom, size)
		assert.EqualValues(t, 15, x)
		assert.EqualValues(t, 10, y)

		// Increasing zoom level by 2 multiplies tile IDs by 4
		x, y = MercatorLocationToTileID(loc.Latitude, loc.Longitude, zoom+2, size)
		assert.EqualValues(t, 15*4+1, x)
		assert.EqualValues(t, 10*4+1, y)

		// Doubling tile size does not change tile ID
		x, y = MercatorLocationToTileID(loc.Latitude, loc.Longitude, zoom, size*2)
		assert.EqualValues(t, 15, x)
		assert.EqualValues(t, 10, y)
	})

	t.Run("Reverses mercator projections", func(t *testing.T) {
		x, y := MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom, size)

		lat2, lng2 := MercatorPixelToLocation(x, y, zoom, size)
		assert.InDelta(t, loc.Latitude, lat2, delta)
		assert.InDelta(t, loc.Longitude, lng2, delta)
	})
}

func BenchmarkMercator(b *testing.B) {
	zoom := uint64(4)
	loc := base.Location{Latitude: -45.942805, Longitude: 166.568500}
	size := uint64(256)

	x, y := MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom, size)

	b.Run("Forward projection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			MercatorLocationToPixel(loc.Latitude, loc.Longitude, zoom, size)
		}
	})

	b.Run("Reverse projection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			MercatorPixelToLocation(x, y, zoom, size)
		}
	})

}
