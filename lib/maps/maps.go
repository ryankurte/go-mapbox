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
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/url"
	"strings"
)

import (
	"github.com/ryankurte/go-mapbox/lib/base"
	"log"
)

const (
	apiName    = "maps"
	apiVersion = "v4"
)

// Cache interface defines an abstract tile cache
// This can be used to limit the number of API calls required to fetch previously fetched tiles
type Cache interface {
	Save(mapID MapID, x, y, level uint64, format MapFormat, highDPI bool, img image.Image) error
	Fetch(mapID MapID, x, y, level uint64, format MapFormat, highDPI bool) (image.Image, *image.Config, error)
}

// Maps api wrapper instance
type Maps struct {
	base  *base.Base
	cache Cache
}

// NewMaps Create a new Maps API wrapper
func NewMaps(base *base.Base) *Maps {
	return &Maps{base, nil}
}

// SetCache binds a cache into the map instance
func (m *Maps) SetCache(c Cache) {
	m.cache = c
}

// GetTile fetches the map tile for the specified location
func (m *Maps) GetTile(mapID MapID, x, y, z uint64, format MapFormat, highDPI bool) (image.Image, *image.Config, error) {

	v := url.Values{}

	// Catch invalid MapID / MapFormat combinations here
	if mapID == MapIDSatellite && strings.Contains(string(format), "png") {
		return nil, nil, fmt.Errorf("MapIDSatellite does not support png outputs")
	}
	if format == MapFormatPngRaw && mapID != MapIDTerrainRGB {
		return nil, nil, fmt.Errorf("MapFormatPngRaw only supported for MapIDTerrainRGB")
	}

	// Attempt cache lookup if available
	if m.cache != nil {
		img, cfg, err := m.cache.Fetch(mapID, x, y, z, format, highDPI)
		if err != nil {
			log.Printf("Cache fetch error (%s)", err)
		} else if img != nil {
			return img, cfg, nil
		}
	}

	// Create Request
	dpiFlag := ""
	if highDPI {
		dpiFlag = "@2x"
	}
	queryString := fmt.Sprintf("%s/%s/%d/%d/%d%s.%s", apiVersion, mapID, z, x, y, dpiFlag, format)

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
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	// Convert to image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	// Save to cache if available
	if m.cache != nil {
		err = m.cache.Save(mapID, x, y, z, format, highDPI, img)
		if err != nil {
			log.Printf("Cache save error (%s)", err)
		}
	}

	return img, &cfg, err
}

// GetEnclosingTiles fetches a 2d array of the tiles enclosing a given point
func (m *Maps) GetEnclosingTiles(mapID MapID, a, b base.Location, level uint64, format MapFormat, highDPI bool) ([][]image.Image, [][]image.Config, error) {
	// Convert to tile locations
	xStart, yStart, xEnd, yEnd := GetEnclosingTileIDs(a, b, level)
	xLen := xEnd - xStart + 1
	yLen := yEnd - yStart + 1

	images := make([][]image.Image, yLen)
	configs := make([][]image.Config, yLen)

	for y := int64(0); y < yLen; y++ {
		images[y] = make([]image.Image, xLen)
		configs[y] = make([]image.Config, xLen)

		for x := int64(0); x < xLen; x++ {

			xIndex := uint64(xStart + x)
			yIndex := uint64(yStart + y)

			xIndex, yIndex = WrapTileID(xIndex, yIndex, level)

			img, cfg, err := m.GetTile(mapID, xIndex, yIndex, level, format, highDPI)
			if err != nil {
				return nil, nil, err
			}

			images[y][x] = img
			configs[y][x] = *cfg
		}
	}

	return images, configs, nil
}
