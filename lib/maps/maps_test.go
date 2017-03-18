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
	//b.SetDebug(true)

	maps := NewMaps(b)

	t.Run("Can fetch map tiles as png", func(t *testing.T) {

		img, _, err := maps.GetTile(MapIDStreets, 1, 0, 1, MapFormatPng, true)

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

		img, _, err := maps.GetTile(MapIDSatellite, 1, 0, 1, MapFormatJpg90, true)

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

		img, _, err := maps.GetTile(MapIDTerrainRGB, 1, 0, 1, MapFormatPngRaw, true)
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

	t.Run("Can fetch map tiles by location", func(t *testing.T) {

		locA := base.Location{-122.42, 20.78}
		locB := base.Location{-77.03, 38.91}

		images, configs, err := maps.GetEnclosingTiles(MapIDSatellite, locA, locB, 4, MapFormatJpg90, true)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		img := maps.StitchTiles(images, configs)

		f, err := os.Create("/tmp/go-mapbox-stitch.jpg")
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
