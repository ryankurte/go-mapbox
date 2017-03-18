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
	"bufio"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

import (
	"github.com/ryankurte/go-mapbox/lib/base"
)

func TestMaps(t *testing.T) {

	token := os.Getenv("MAPBOX_TOKEN")
	if token == "" {
		t.Error("Mapbox API token not found")
		t.FailNow()
	}

	b := base.NewBase(token)

	maps := NewMaps(b)

	t.Run("Can fetch map tiles as png", func(t *testing.T) {

		img, err := maps.GetTile(MapIDStreets, 1, 0, 0, MapFormatPng, true)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		f, err := os.Create("/tmp/go-mapbox-test.png")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		w := bufio.NewWriter(f)

		err = png.Encode(w, img)
		if err != nil {
			t.Error(err)
		}

		f.Close()
	})

	t.Run("Can fetch map tiles as jpeg", func(t *testing.T) {

		img, err := maps.GetTile(MapIDSatellite, 1, 0, 0, MapFormatJpg90, true)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		f, err := os.Create("/tmp/go-mapbox-test.jpg")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		w := bufio.NewWriter(f)

		err = jpeg.Encode(w, img, nil)
		if err != nil {
			t.Error(err)
		}

		f.Close()
	})

	t.Run("Can fetch terrain RGB tiles", func(t *testing.T) {

		img, err := maps.GetTile(MapIDTerrainRGB, 1, 0, 0, MapFormatPngRaw, true)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		f, err := os.Create("/tmp/go-mapbox-test-terrain.png")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		w := bufio.NewWriter(f)

		err = jpeg.Encode(w, img, nil)
		if err != nil {
			t.Error(err)
		}

		f.Close()
	})

}
