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
	"io/ioutil"
	"os"
)

import (
	"bytes"
	"github.com/paulmach/go.geo"
	"github.com/ryankurte/go-mapbox/lib/base"
)

// LocationToTileID converts a lat/lon location into a tile ID
func LocationToTileID(loc base.Location, level uint64) (int64, int64) {
	x, y := geo.ScalarMercator.Project(loc.Longitude, loc.Latitude, level)
	return int64(x), int64(y)
}

// TileIDToLocation converts a tile ID to a lat/lon location
func TileIDToLocation(x, y, level uint64) base.Location {
	lat, lng := geo.ScalarMercator.Inverse(x, y, level)
	return base.Location{
		Latitude:  lat,
		Longitude: lng,
	}
}

// WrapTileID wraps tile IDs by level for api requests
// eg. Tile (X:16, Y:10, level:4 ) will become (X:0, Y:10, level:4)
func WrapTileID(x, y, level uint64) (uint64, uint64) {
	// Limit to 2^n tile range for a given level
	x = x % (2 << (level - 1))
	y = y % (2 << (level - 1))

	return x, y
}

// GetEnclosingTileIDs fetches a pair of tile IDs enclosing the provided pair of points
func GetEnclosingTileIDs(a, b base.Location, level uint64) (int64, int64, int64, int64) {
	aX, aY := LocationToTileID(a, level)
	bX, bY := LocationToTileID(b, level)

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
func StitchTiles(images [][]image.Image, config image.Config) image.Image {

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

// LoadImage loads an image from a file
func LoadImage(file string) (image.Image, *image.Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}

	r := bufio.NewReader(f)
	data, err := ioutil.ReadAll(r)
	f.Close()

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		f.Close()
		return nil, nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		f.Close()
		return nil, nil, err
	}

	return img, &cfg, nil
}

// SaveImageJPG writes an image instance to a jpg file
func SaveImageJPG(img image.Image, file string) error {
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
func SaveImagePNG(img image.Image, file string) error {
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

// PixelToHeight Converts a pixel to a height value for mapbox terrain tiles
// Equation from https://www.mapbox.com/blog/terrain-rgb/
func PixelToHeight(pixel image.RGBA) float64 {
	R := float64(pixel.Pix[0])
	G := float64(pixel.Pix[1])
	B := float64(pixel.Pix[2])
	return -10000 + ((R*256*256 + G*256 + B) * 0.1)
}
