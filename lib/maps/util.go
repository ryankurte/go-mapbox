/**
 * go-mapbox Maps Module Utils
 * Utilities and Helpers for dealing with maps and map tiles
 * See https://www.mapbox.com/api-documentation/#maps for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package maps

import (
	"bufio"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

import (
	"github.com/paulmach/go.geo"
	"github.com/ryankurte/go-mapbox/lib/base"
)

// LocationToTileID converts a lat/lon location into a tile ID
func (m *Maps) LocationToTileID(loc base.Location, level uint64) (int64, int64) {
	if loc.Latitude > 180.0 {
		loc.Latitude = 360.0 - loc.Latitude
	}

	// Calculate X and Y at a given level
	x, y := geo.ScalarMercator.Project(loc.Longitude, loc.Latitude, level)

	log.Printf("Tile location (%d, %d)", x, y)

	log.Printf("Corrected tile location (%d, %d) limit %d", x, y, (2 << (level - 1)))

	return int64(x), int64(y)
}

// TileIDToLocation converts a tile ID to a lat/lon location
func (m *Maps) TileIDToLocation(x, y, level uint64) base.Location {

	lat, lng := geo.ScalarMercator.Inverse(x, y, level)

	return base.Location{
		Latitude:  lat,
		Longitude: lng,
	}
}

// WrapTileID wraps tile IDs by level for api requests
// eg. Tile (X:16, Y:10, level:4 )will become (X:0, Y:10, level:4)
func (m *Maps) WrapTileID(x, y, level uint64) (uint64, uint64) {
	// Limit to 2^n tile range for a given level
	x = x % (2 << (level - 1))
	y = y % (2 << (level - 1))

	return x, y
}

// GetEnclosingTileIDs fetches a pair of tile IDs enclosing the provided pair of points
func (m *Maps) GetEnclosingTileIDs(a, b base.Location, level uint64) (int64, int64, int64, int64) {
	aX, aY := m.LocationToTileID(a, level)
	bX, bY := m.LocationToTileID(b, level)

	log.Printf("aX: %d aY: %d bX: %d bY: %d", aX, aY, bX, bY)

	var xStart, xEnd, yStart, yEnd int64
	if bX >= aX {
		xStart = aX
		xEnd = bX
	} else {
		xStart = bX
		xEnd = aX
	}

	if bY >= aY {
		yStart = aY
		yEnd = bY
	} else {
		yStart = bY
		yEnd = aY
	}

	return xStart, yStart, xEnd, yEnd
}

// StitchTiles combines a 2d array of image tiles into a single larger image
// Note that all images must have the same dimensions for this to work
func (m *Maps) StitchTiles(images [][]image.Image, config image.Config) image.Image {

	imgX := config.Width
	imgY := config.Height

	xSize := imgX * len(images[0])
	ySize := imgY * len(images)

	stitched := image.NewRGBA(image.Rect(0, 0, xSize, ySize))

	for y, row := range images {
		for x, img := range row {
			sp := image.Point{0, 0}
			bounds := image.Rect(x*imgX, y*imgY, (x+1)*imgX, (y+1)*imgY)
			draw.Draw(stitched, bounds, img, sp, draw.Over)
		}
	}

	return stitched
}

// SaveImageJPG writes an image instance to a jpg file
func (m *Maps) SaveImageJPG(img image.Image, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)

	err = jpeg.Encode(w, img, nil)
	if err != nil {
		return err
	}

	f.Close()

	return nil
}

// SaveImagePNG writes an image instance to a png file
func (m *Maps) SaveImagePNG(img image.Image, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)

	err = png.Encode(w, img)
	if err != nil {
		return err
	}

	f.Close()

	return nil
}
