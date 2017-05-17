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

func TestScalarMercator(t *testing.T) {

	t.Run("Performs naive web mercator projections", func(t *testing.T) {
		zoom := uint64(4)
		size := uint64(256)
		fsize := float64(size)

		loc := base.Location{-45.942805, 166.568500}

		xExpected, yExpected := LocationToTileID(loc, zoom)

		x, y := NaiveLocationToPixel(loc.Latitude, loc.Longitude, zoom, size)
		assert.EqualValues(t, xExpected, math.Floor(x/fsize))
		assert.EqualValues(t, yExpected, math.Floor(y/fsize))

		lat2, lng2 := NaivePixelToLocation(x, y, zoom, size)
		assert.InDelta(t, loc.Latitude, lat2, delta)
		assert.InDelta(t, loc.Longitude, lng2, delta)
	})

	t.Run("Performs cached scalar mercator projections", func(t *testing.T) {

		zoom := uint64(4)
		loc := base.Location{-45.942805, 166.568500}

		sm := NewSphericalMercator(256)

		xExpected, yExpected := LocationToTileID(loc, zoom)

		x, y := sm.LocationToPixel(loc.Latitude, loc.Longitude, zoom)
		assert.EqualValues(t, xExpected, math.Floor(x/256))
		assert.EqualValues(t, yExpected, math.Floor(y/256))

		lat2, lng2 := sm.PixelToLocation(x, y, zoom)
		assert.InDelta(t, loc.Latitude, lat2, delta)
		assert.InDelta(t, loc.Longitude, lng2, delta)
	})

}

func BenchmarkScalarMercator(b *testing.B) {
	zoom := uint64(4)
	loc := base.Location{Latitude: -45.942805, Longitude: 166.568500}
	size := uint64(256)

	sm := NewSphericalMercator(256)
	x, y := NaiveLocationToPixel(loc.Latitude, loc.Longitude, zoom, size)

	b.Run("Forward naive projection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NaiveLocationToPixel(loc.Latitude, loc.Longitude, zoom, size)
		}
	})

	b.Run("Reverse naive projection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NaiveLocationToPixel(x, y, zoom, size)
		}
	})

	b.Run("Forward optimised projection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sm.LocationToPixel(loc.Latitude, loc.Longitude, zoom)
		}
	})

	b.Run("Reverse optimised projection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sm.LocationToPixel(x, y, zoom)
		}
	})

}
