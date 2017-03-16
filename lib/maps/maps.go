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
	"github.com/ryankurte/go-mapbox/lib/base"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/url"
	"strings"
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
func (m *Maps) GetTiles(mapID MapID, x, y, z uint, format MapFormat, highDPI bool) (image.Image, error) {

	v := url.Values{}

	// Catch invalid MapID / MapFormat combinations here
	if mapID == MapIDSatellite && strings.Contains(string(format), "png") {
		return nil, fmt.Errorf("MapIDSatellite does not support png outputs")
	}

	// Create Request
	dpiFlag := ""
	if highDPI {
		dpiFlag = "@2x"
	}
	queryString := fmt.Sprintf("%s/%s/%d/%d/%d%s.%s", apiVersion, mapID, x, y, z, dpiFlag, format)
	resp, err := m.base.QueryRequest(queryString, &v)

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
	defer resp.Body.Close()

	// Convert to image
	var img image.Image
	reader := bytes.NewReader(data)
	switch contentType {
	case "image/png":
		img, err = png.Decode(reader)
	case "image/jpg":
		img, err = jpeg.Decode(reader)
	case "image/jpeg":
		img, err = jpeg.Decode(reader)
	default:
		return nil, fmt.Errorf("Unrecognised Content-Type (%s)", contentType)
	}

	return img, err
}
