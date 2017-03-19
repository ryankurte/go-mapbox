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
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"os"
)

import (
	"github.com/ryankurte/go-mapbox/lib/base"
)

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

// Mercator transformations for calculating map positions
// http://wiki.openstreetmap.org/wiki/Slippy_map_tilenames

// LatLonToTileXY converts a lat and lon (in degrees) at a given zoom level to tile coordinates
func LatLonToTileXY(lat, lon float64, zoom uint64) (float64, float64) {
	n := math.Pow(2, float64(zoom))

	latRads := lat / 180.0 * math.Pi

	x := n * ((lon + 180.0) / 360.0)
	y := n * (1 - (math.Log(math.Tan(latRads)+1/math.Cos(latRads)) / math.Pi)) / 2

	return x, y
}

// TileXYToLatLon converts a tile position (x, y) at a given zoom level to a lat and lon
func TileXYToLatLon(x, y float64, zoom uint64) (float64, float64) {
	n := math.Pow(2, float64(zoom))

	lonDeg := (x / n * 360.0) - 180.0
	latRad := math.Atan(math.Sinh(math.Pi * (1 - (2 * y / n))))
	latDeg := latRad * 180.0 / math.Pi

	return latDeg, lonDeg
}

// LocationToTileID converts a lat/lon location into a tile ID
func LocationToTileID(loc base.Location, level uint64) (int64, int64) {
	x, y := LatLonToTileXY(loc.Latitude, loc.Longitude, level)
	return int64(math.Floor(x)), int64(math.Floor(y))
}
