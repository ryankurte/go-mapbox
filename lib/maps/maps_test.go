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
	"os"
	"testing"
)

import (
	"fmt"
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

		err = SaveImagePNG(img, "/tmp/go-mapbox-test.png")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	})

	t.Run("Can fetch map tiles as jpeg", func(t *testing.T) {

		img, _, err := maps.GetTile(MapIDSatellite, 1, 0, 1, MapFormatJpg90, true)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		err = SaveImageJPG(img, "/tmp/go-mapbox-test.jpg")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	})

	t.Run("Can fetch terrain RGB tiles", func(t *testing.T) {

		img, _, err := maps.GetTile(MapIDTerrainRGB, 1, 0, 1, MapFormatPngRaw, true)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		err = SaveImagePNG(img, "/tmp/go-mapbox-terrain.png")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	})

	t.Run("Can fetch map tiles by location", func(t *testing.T) {

		locA := base.Location{-45.942805, 166.568500}
		locB := base.Location{-34.2186101, 183.4015517}

		images, configs, err := maps.GetEnclosingTiles(MapIDSatellite, locA, locB, 6, MapFormatJpg90, true)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		for y := range images {
			for x := range images[y] {
				SaveImageJPG(images[y][x], fmt.Sprintf("/tmp/go-mapbox-stitch-%d-%d.jpg", x, y))
			}
		}

		img := StitchTiles(images, configs[0][0])

		err = SaveImageJPG(img, "/tmp/go-mapbox-stitch.jpg")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

	})

	t.Run("Can fetch map tiles by location (with cache)", func(t *testing.T) {

		locA := base.Location{-45.942805, 166.568500}
		locB := base.Location{-34.2186101, 183.4015517}

		cache, err := NewFileCache("/tmp/go-mapbox-cache")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		maps.SetCache(cache)

		images, configs, err := maps.GetEnclosingTiles(MapIDSatellite, locA, locB, 6, MapFormatJpg90, true)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		img := StitchTiles(images, configs[0][0])

		err = SaveImageJPG(img, "/tmp/go-mapbox-stitch2.jpg")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

	})

}
