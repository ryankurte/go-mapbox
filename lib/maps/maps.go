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
	"log"
	"net/url"
	"strings"

	"github.com/JayBusch/go-mapbox/lib/base"
	"sync"
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
func (m *Maps) GetTile(mapID MapID, x, y, z uint64, format MapFormat, highDPI bool) (*Tile, error) {

	v := url.Values{}

	dpiFlag := ""
	size := SizeStandard
	if highDPI {
		dpiFlag = "@2x"
		size = SizeHighDPI
	}

	// Catch invalid MapID / MapFormat combinations here
	if mapID == MapIDSatellite && strings.Contains(string(format), "png") {
		return nil, fmt.Errorf("MapIDSatellite does not support png outputs")
	}
	if format == MapFormatPngRaw && mapID != MapIDTerrainRGB {
		return nil, fmt.Errorf("MapFormatPngRaw only supported for MapIDTerrainRGB")
	}
	if mapID == MapIDTerrainRGB && format != MapFormatPngRaw {
		return nil, fmt.Errorf("MapIDTerrainRGB only supports format MapFormatPngRaw")
	}

	// Attempt cache lookup if available
	if m.cache != nil {
		img, _, err := m.cache.Fetch(mapID, x, y, z, format, highDPI)
		if err != nil {
			log.Printf("Cache fetch error (%s)", err)
		} else if img != nil {
			tile := NewTile(x, y, z, size, img)
			return &tile, nil
		}
	}

	// Create Request
	queryString := fmt.Sprintf("%s/%s/%d/%d/%d%s.%s", apiVersion, mapID, z, x, y, dpiFlag, format)

	resp, err := m.base.QueryRequest(queryString, &v)
	if err != nil {
		return nil, err
	}

	// Parse content type and length
	contentType := resp.Header.Get("Content-Type")
	contentLength := resp.ContentLength

	// Read data from body
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body (%s)", err)
	}
	if len(data) != int(contentLength) {
		return nil, fmt.Errorf("Content length mismatch (expected %d received %d)", contentLength, len(data))
	}

	if strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("Invalid API call: %s message: %s", resp.Request.URL, string(data))
	}

	// Decode config
	_, _, err = image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// Convert to image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// Create tile
	tile := NewTile(x, y, z, size, img)

	// Save to cache if available
	// Tile is post RGB conversion (should avoid pngraw issues)
	if m.cache != nil {
		err = m.cache.Save(mapID, x, y, z, format, highDPI, img)
		if err != nil {
			log.Printf("Cache save error (%s)", err)
		}
	}

	return &tile, err
}

// GetEnclosingTiles fetches a 2d array of the tiles enclosing a given point
func (m *Maps) GetEnclosingTiles(mapID MapID, a, b base.Location, level uint64, format MapFormat, highDPI bool) ([][]Tile, error) {
	// Convert to tile locations
	xStart, yStart, xEnd, yEnd := GetEnclosingTileIDs(a, b, level)
	xLen := xEnd - xStart + 1
	yLen := yEnd - yStart + 1

	tiles := make([][]Tile, yLen)

	for y := uint64(0); y < yLen; y++ {
		tiles[y] = make([]Tile, xLen)

		for x := uint64(0); x < xLen; x++ {

			xIndex := uint64(xStart + x)
			yIndex := uint64(yStart + y)

			xIndex, yIndex = WrapTileID(xIndex, yIndex, level)

			tile, err := m.GetTile(mapID, xIndex, yIndex, level, format, highDPI)
			if err != nil {
				return nil, err
			}

			tiles[y][x] = *tile
		}
	}

	return tiles, nil
}

func (m *Maps) FastGetEnclosingTiles(mapID MapID, a, b base.Location, level uint64, format MapFormat, highDPI bool) ([][]Tile, error) {
	// Convert to tile locations
	xStart, yStart, xEnd, yEnd := GetEnclosingTileIDs(a, b, level)
	xLen := xEnd - xStart + 1
	yLen := yEnd - yStart + 1

	in := make(chan *Tile, 1)
	var wg1 sync.WaitGroup
	wg1.Add(int(xLen * yLen))

	tiles := make([][]Tile, yLen)
	for y := uint64(0); y < yLen; y++ {
		tiles[y] = make([]Tile, xLen)

		for x := uint64(0); x < xLen; x++ {
			xIndex := uint64(xStart + x)
			yIndex := uint64(yStart + y)

			xIndex, yIndex = WrapTileID(xIndex, yIndex, level)

			go func(xIndex, yIndex uint64) {
				tile, err := m.GetTile(mapID, xIndex, yIndex, level, format, highDPI)
				if err != nil {
					log.Printf("Error fetching tile: %s", err)
				}
				in <- tile
				wg1.Done()
			}(xIndex, yIndex)
		}
	}

	go func() {
		wg1.Wait()
		close(in)
	}()

stitch:
	for {
		select {
		case t, ok := <-in:
			if !ok {
				break stitch
			}
			if t == nil {
				return nil, fmt.Errorf("api error")
			}
			tiles[t.Y-yStart][t.X-xStart] = *t
		}
	}

	return tiles, nil
}
