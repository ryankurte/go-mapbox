/**
 * go-mapbox Maps Module
 * Spherical Mercator implementation
 * Based on https://github.com/mapbox/sphericalmercator/blob/master/sphericalmercator.js
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package maps

import (
	"math"
)

const (
	D2R = math.Pi / 180 // D2R helper for converting degrees to radians
	R2D = 180 / math.Pi // R2D helper for converting radians to degrees

	// 900913 (GOOGLE) spec properties
	A         = 6378137.0
	MAXEXTENT = 20037508.342789244
	d         = 30
)

type cacheLine struct {
	bc []float64
	cc []float64
	zc []float64
	ac []float64
}

var cache map[uint64]cacheLine

func init() {
	cache = make(map[uint64]cacheLine)
}

func buildCache(size uint64) cacheLine {
	cl := cacheLine{
		bc: make([]float64, 30),
		cc: make([]float64, 30),
		zc: make([]float64, 30),
		ac: make([]float64, 30),
	}

	fsize := float64(size)

	for d := 0; d < 30; d++ {
		cl.bc[d] = fsize / 360.0
		cl.cc[d] = fsize / (2 * math.Pi)
		cl.zc[d] = fsize / 2
		cl.ac[d] = fsize
		fsize *= 2
	}
	return cl
}

// SphericalMercator implements spherical mercator transforms for a given tile size
// Tile calculations are pre-cached on first creation for performance, so this can be created / destroyed as required
type SphericalMercator struct {
	size uint64
	bc   []float64
	cc   []float64
	zc   []float64
	ac   []float64
}

// NewSphericalMercator Creates a spherical mercator with the provided tile size
// This uses pre-cached data if available, and adds to the cache if not to decrease
// the runtime cost of creating SphericalMercators
func NewSphericalMercator(size uint64) SphericalMercator {
	sm := SphericalMercator{
		size: size,
	}

	// Create cache line if not already configured
	if _, ok := cache[size]; !ok {
		cache[size] = buildCache(size)
	}

	// Map pre-calculated data from cache
	sm.bc = cache[size].bc
	sm.cc = cache[size].cc
	sm.zc = cache[size].zc
	sm.ac = cache[size].ac

	return sm
}

// LocationToPixel converts a lat/lng/zoom location (in degrees) to a pixel location in the global space
// ie. this is relative to 0, 0 across all tiles.
// The TileID is then this pixel location divided by the tile size as an integer (see LocationToTileID)
func (sm *SphericalMercator) LocationToPixel(lat, lng float64, zoom uint64) (float64, float64) {
	d := sm.zc[zoom]
	f := math.Min(math.Max(math.Sin(D2R*lat), -0.9999), 0.9999)
	x := math.Floor(d + lng*sm.bc[zoom])
	y := math.Floor(d + 0.5*math.Log((1+f)/(1-f))*(-sm.cc[zoom]))
	if x > sm.ac[zoom] {
		x = sm.ac[zoom]
	}
	if y > sm.ac[zoom] {
		y = sm.ac[zoom]
	}
	return x, y
}

// LocationToTileID builds on LocationToPixel to fetch the TileID of a given location
func (sm *SphericalMercator) LocationToTileID(lat, lng float64, zoom uint64) (uint64, uint64) {
	size := float64(sm.size)
	x, y := sm.LocationToPixel(lat, lng, zoom)
	return uint64(x / size), uint64(y / size)
}

// PixelToLocation converts a given (global) pixel location and zoom level to a lat and lng (in degrees)
func (sm *SphericalMercator) PixelToLocation(x, y float64, zoom uint64) (float64, float64) {
	g := (y - sm.zc[zoom]) / (-sm.cc[zoom])
	lng := (x - sm.zc[zoom]) / sm.bc[zoom]
	lat := R2D * (2*math.Atan(math.Exp(g)) - 0.5*math.Pi)
	return lat, lng
}

// NaiveLocationToPixel converts a lat/lng/zoom location (in degrees) to a pixel location in the global space
func NaiveLocationToPixel(lat, lng float64, zoom, size uint64) (float64, float64) {
	pi := math.Pi
	latRad, lngRad := lat*D2R, lng*D2R

	//fsize := float64(size)
	x := (128 / pi) * math.Pow(2, float64(zoom)) * (lngRad + pi)
	y := (128 / pi) * math.Pow(2, float64(zoom)) * (pi - math.Log(math.Tan(pi/4+latRad/2)))
	return x, y
}

// NaiveLocationToTileID builds on NaiveLocationToPixel to fetch the TileID of a given location
func NaiveLocationToTileID(lat, lng float64, zoom, size uint64) (uint64, uint64) {
	fsize := float64(size)
	x, y := NaiveLocationToPixel(lat, lng, zoom, size)
	xID, yID := uint64(x/fsize), uint64(y/fsize)
	return xID, yID
}

// NaivePixelToLocation converts a given (global) pixel location and zoom level to a lat and lng (in degrees)
func NaivePixelToLocation(x, y float64, zoom, size uint64) (float64, float64) {
	pi := math.Pi
	fsize := float64(size)
	lng := x*(pi/fsize*2)/math.Pow(2, float64(zoom)) - pi
	lat := (math.Atan(math.Pow(math.E, (-y*(pi/fsize*2)/math.Pow(2, float64(zoom))+pi))) - pi/4) * 2
	return lat * R2D, lng * R2D
}
