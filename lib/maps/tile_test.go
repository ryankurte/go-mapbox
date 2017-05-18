/**
 * go-mapbox Maps Module Tile HelperTests
 * See https://www.mapbox.com/api-documentation/#maps for API information
 * Run: go test -v  -run ^TestTiles$ && open /tmp/mapbox-tile-test-*.jpg
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package maps

import (
	"image/color"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryankurte/go-mapbox/lib/base"
)

func TestTiles(t *testing.T) {

	token := os.Getenv("MAPBOX_TOKEN")
	if token == "" {
		t.Error("Mapbox API token not found")
		t.FailNow()
	}

	b := base.NewBase(token)

	maps := NewMaps(b)

	cache, err := NewFileCache("/tmp/go-mapbox-cache")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	maps.SetCache(cache)

	size := uint64(512)
	x, y, z := uint64(15), uint64(9), uint64(4)

	loc := base.Location{-36.8485, 174.7633}

	img, _, err := maps.GetTile(MapIDSatellite, x, y, z, MapFormatJpg90, true)
	assert.Nil(t, err)

	err = SaveImageJPG(img, "/tmp/mapbox-tile-test-1.jpg")
	assert.Nil(t, err)

	fire, _, err := LoadImage("../../fire64.png")
	assert.Nil(t, err)

	t.Run("Can draw in local tile space", func(t *testing.T) {
		tile := NewTile(x, y, z, size, img)
		tile.DrawLocalXY(fire, int(size/2), int(size/2), Center)

		err := SaveImageJPG(tile, "/tmp/mapbox-tile-test-2.jpg")
		assert.Nil(t, err)
	})

	t.Run("Can draw in global tile space", func(t *testing.T) {
		tile := NewTile(x, y, z, size, img)
		err := tile.DrawGlobalXY(fire, int(size*x+size/2), int(size*y+size/2), Center)
		assert.Nil(t, err)

		err = SaveImageJPG(tile, "/tmp/mapbox-tile-test-3.jpg")
		assert.Nil(t, err)
	})

	t.Run("Can draw in world space", func(t *testing.T) {
		tile := NewTile(x, y, z, size, img)

		tile.DrawLocation(fire, loc, DrawConfig{Vertical: JustifyBottom, Horizontal: JustifyCenter})
		tile.DrawLocation(fire, loc, DrawConfig{Vertical: JustifyCenter, Horizontal: JustifyLeft})
		tile.DrawLocation(fire, loc, DrawConfig{Vertical: JustifyCenter, Horizontal: JustifyCenter})
		tile.DrawLocation(fire, loc, DrawConfig{Vertical: JustifyCenter, Horizontal: JustifyRight})
		tile.DrawLocation(fire, loc, DrawConfig{Vertical: JustifyTop, Horizontal: JustifyCenter})

		err := SaveImageJPG(tile, "/tmp/mapbox-tile-test-4.jpg")
		assert.Nil(t, err)
	})

	t.Run("Can render to composite tiles", func(t *testing.T) {
		locA := base.Location{-45.942805, 166.568500}
		locB := base.Location{-34.2186101, 183.4015517}

		x1, y1, _, _ := GetEnclosingTileIDs(locA, locB, 6)
		images, configs, err := maps.GetEnclosingTiles(MapIDSatellite, locA, locB, 6, MapFormatJpg90, true)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		img := StitchTiles(images, configs[0][0])

		tile := NewTile(x1, y1, 6, size, img)

		tile.DrawLocation(fire, base.Location{-41.2865, 174.7762}, DrawConfig{Vertical: JustifyBottom, Horizontal: JustifyCenter})
		tile.DrawLocation(fire, base.Location{-36.8485, 174.7633}, DrawConfig{Vertical: JustifyBottom, Horizontal: JustifyCenter})
		tile.DrawLocation(fire, base.Location{-43.5321, 172.6362}, DrawConfig{Vertical: JustifyBottom, Horizontal: JustifyCenter})

		err = SaveImageJPG(tile, "/tmp/mapbox-tile-test-5.jpg")
		assert.Nil(t, err)
	})

	t.Run("Can interpolate lines over complex tiles", func(t *testing.T) {
		locA := base.Location{-45.942805, 166.568500}
		locB := base.Location{-34.2186101, 183.4015517}

		x1, y1, _, _ := GetEnclosingTileIDs(locA, locB, 6)
		images, configs, err := maps.GetEnclosingTiles(MapIDSatellite, locA, locB, 6, MapFormatJpg90, true)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		img := StitchTiles(images, configs[0][0])

		tile := NewTile(x1, y1, 6, size, img)
		a, b, c := base.Location{-36.8485, 174.7633}, base.Location{-41.2865, 174.7762}, base.Location{-43.5321, 172.6362}
		tile.DrawLine(a, b, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		tile.DrawLine(b, c, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		tile.DrawLine(c, a, color.RGBA{R: 255, G: 0, B: 0, A: 255})

		err = SaveImageJPG(tile, "/tmp/mapbox-tile-test-6.jpg")
		assert.Nil(t, err)
	})

	t.Run("Can fetch terrain data points", func(t *testing.T) {
		locA := base.Location{-39.5, 173.5}
		locB := base.Location{-39.0, 174.5}
		taranaki := base.Location{-39.295182, 174.063668}
		level := uint64(11)

		images, configs, err := maps.GetEnclosingTiles(MapIDTerrainRGB, locA, locB, level, MapFormatPngRaw, true)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		img := StitchTiles(images, configs[0][0])

		x1, y1, _, _ := GetEnclosingTileIDs(locA, locB, level)
		tile := NewTile(x1, y1, level, size, img)

		err = SaveImageJPG(tile, "/tmp/mapbox-tile-test-7.png")
		assert.Nil(t, err)

		alt, err := tile.GetAltitude(taranaki)
		assert.Nil(t, err)
		assert.InDelta(t, 2400, alt, 100)

		flattened := tile.FlattenAltitudes(3000)
		err = SaveImageJPG(flattened, "/tmp/mapbox-tile-test-8.png")
		assert.Nil(t, err)

	})

}
