/**
 * go-mapbox Maps Module
 * Wraps the mapbox geocoding API for server side use
 * See https://www.mapbox.com/api-documentation/#maps for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package maps

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"image/draw"
	"io/ioutil"
	"net/url"
	"strings"
	"log"
)

import (
	"github.com/ryankurte/go-mapbox/lib/base"
	"github.com/paulmach/go.geo"

)

const (
	apiName    = "maps"
	apiVersion = "v4"
)

// Maps api wrapper instance
type Maps struct {
	base *base.Base
}

// NewMaps Create a new Maps API wrapper
func NewMaps(base *base.Base) *Maps {
	return &Maps{base}
}

// GetTiles fetches the map tile for the specified location
func (m *Maps) GetTiles(mapID MapID, x, y, z uint64, format MapFormat, highDPI bool) (image.Image, *image.Config, error) {

	v := url.Values{}

	// Catch invalid MapID / MapFormat combinations here
	if mapID == MapIDSatellite && strings.Contains(string(format), "png") {
		return nil, nil, fmt.Errorf("MapIDSatellite does not support png outputs")
	}

	// Create Request
	dpiFlag := ""
	if highDPI {
		dpiFlag = "@2x"
	}
	queryString := fmt.Sprintf("%s/%s/%d/%d/%d%s.%s", apiVersion, mapID, z, x, y, dpiFlag, format)
	log.Printf("Fetching tile (x: %d, y: %d, z: %d) query: %s", x, y, z, queryString)

	resp, err := m.base.QueryRequest(queryString, &v)
	if err != nil {
		return nil, nil, err
	}


	// Parse content type and length
	contentType := resp.Header.Get("Content-Type")
	contentLength := resp.ContentLength

	// Read data from body
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("Error reading response body (%s)", err)
	}
	if len(data) != int(contentLength) {
		return nil, nil, fmt.Errorf("Content length mismatch (expected %d received %d)", contentLength, len(data))
	}

	if strings.Contains(contentType, "application/json") {
		return nil, nil, fmt.Errorf("Invalid API call: %s message: %s", resp.Request.URL, string(data))
	}

	// Load config
	reader := bytes.NewReader(data)
	var cfg image.Config
	err = nil
	switch contentType {
	case "image/png":
		cfg, err = png.DecodeConfig(reader)
	case "image/jpg":
		cfg, err = jpeg.DecodeConfig(reader)
	case "image/jpeg":
		cfg, err = jpeg.DecodeConfig(reader)
	default:
		return nil, nil, fmt.Errorf("Unrecognised Content-Type (%s)", contentType)
	}

	if err != nil {
		return nil, nil, err
	}

	// Convert to image
	var img image.Image
	reader = bytes.NewReader(data)
	switch contentType {
	case "image/png":
		img, err = png.Decode(reader)
	case "image/jpg":
		img, err = jpeg.Decode(reader)
	case "image/jpeg":
		img, err = jpeg.Decode(reader)
	default:
		return nil, nil, fmt.Errorf("Unrecognised Content-Type (%s)", contentType)
	}

	return img, &cfg, err
}

func (m *Maps) GetEnclosingTiles(mapID MapID, a, b base.Location, level uint64, format MapFormat, highDPI bool) ([][]image.Image, [][]image.Config, error) {
	// Convert to tile locations
	aX, aY := geo.ScalarMercator.Project(a.Longitude, a.Latitude, level)
	bX, bY := geo.ScalarMercator.Project(b.Longitude, b.Latitude, level)

	log.Printf("aX: %d aY: %d bX: %d bY: %d", aX, aY, bX, bY)

	var xStart, xEnd, yStart, yEnd int64
	if bX >= aX {
		xStart = int64(aX)
		xEnd = int64(bX)
	} else {
		xStart = int64(bX)
		xEnd = int64(aX)
	}
	xLen := xEnd - xStart

	if bY >= aY {
		yStart = int64(aY)
		yEnd = int64(bY)
	} else {
		yStart = int64(bY)
		yEnd = int64(aY)
	}
	yLen := yEnd - yStart

	log.Printf("X (start: %d end: %d len: %d) Y (start: %d end: %d len: %d)", xStart, xEnd, xLen, yStart, yEnd, yLen)

	log.Printf("Fetching %d x %d tiles from (%d, %d) to (%d, %d) at level %x", xLen, yLen, xStart, yStart, xEnd, yEnd, level)

	images := make([][]image.Image, yLen)
	configs := make([][]image.Config, yLen)

	log.Printf("Images: %+v", images)

	count := 0
	for y := int64(0); y < yLen ; y += 1 {
		images[y] = make([]image.Image, xLen)
		configs[y] = make([]image.Config, xLen)

		for x := int64(0); x < xLen; x += 1  {

			xIndex := uint64(xStart + x)
			yIndex := uint64(yStart + y)

			log.Printf("Iteration %d Fetching tile (x: %d, y: %d, z: %d)", count, xIndex, yIndex, level)

			img, cfg, err := m.GetTiles(mapID, xIndex, yIndex, level, format, highDPI)
			if err != nil {
				return nil, nil, err
			}

			images[y][x] = img
			configs[y][x] = *cfg

			count ++
		}
	} 

	return images, configs, nil
}

func (m *Maps) StitchTiles(images [][]image.Image, configs [][]image.Config) image.Image {

	imgX := configs[0][0].Width
	imgY := configs[0][0].Height

	xSize := imgX * len(images[0]);
	ySize := imgY * len(images);


	stitched := image.NewRGBA(image.Rect(0, 0, xSize, ySize))


	for y, row := range images {
		for x, img := range row {

			sp := image.Point{x * imgX, y * imgY}
			draw.Draw(stitched, img.Bounds(), img, sp, draw.Over)
		}
	}


	return stitched
}
