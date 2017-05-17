package maps

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/ryankurte/go-mapbox/lib/base"
)

// Tile is a wrapper around an image that includes positioning data
type Tile struct {
	image.Image
	Level uint64 // Tile zoom level
	Size  uint64 // Tile size
	X, Y  uint64 // Tile X and Y postions (Web Mercurator projection)
}

// NewTile creates a tile with a base RGBA object
func NewTile(x, y, level, size uint64, src image.Image) Tile {
	// Convert image to RGBA
	b := src.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), src, b.Min, draw.Src)

	return Tile{
		Image: m,
		X:     x,
		Y:     y,
		Level: level,
		Size:  size,
	}
}

// Justify sets image offsets for drawing
type Justify string

const (
	JustifyTop    Justify = "top"
	JustifyLeft   Justify = "left"
	JustifyCenter Justify = "center"
	JustifyBottom Justify = "bottom"
	JustifyRight  Justify = "right"
)

type DrawConfig struct {
	Vertical   Justify
	Horizontal Justify
}

var Center = DrawConfig{JustifyCenter, JustifyCenter}

// DrawLocalXY draws the provided image at the local X/Y coordinates
func (t *Tile) DrawLocalXY(src image.Image, x, y int, config DrawConfig) error {
	dp := image.Point{}

	switch config.Horizontal {
	case JustifyLeft:
		dp.X = x
	case JustifyCenter:
		dp.X = x - src.Bounds().Dx()/2
	case JustifyRight:
		dp.X = x - src.Bounds().Dx()
	default:
		return fmt.Errorf("Unsupported horizontal justification (%s)", config.Horizontal)
	}

	switch config.Vertical {
	case JustifyTop:
		dp.Y = y
	case JustifyCenter:
		dp.Y = y - src.Bounds().Dy()/2
	case JustifyBottom:
		dp.Y = y - src.Bounds().Dy()
	default:
		return fmt.Errorf("Unsupported vertical justification (%s)", config.Horizontal)
	}

	r := image.Rectangle{dp, dp.Add(src.Bounds().Size())}
	draw.Draw(t.Image.(draw.Image), r, src, src.Bounds().Min, draw.Over)

	return nil
}

// DrawGlobalXY draws the provided image at the global X/Y coordinates
func (t *Tile) DrawGlobalXY(src image.Image, x, y int, config DrawConfig) error {
	offsetX := x - int(t.X*t.Size)
	offsetY := y - int(t.Y*t.Size)

	if (offsetX < 0) || (offsetX > int(t.Image.Bounds().Max.X)) {
		return fmt.Errorf("Tile DrawGlobalXY error: global X offset not within tile space (%d)", offsetX)
	}
	if (offsetY < 0) || (offsetY > int(t.Image.Bounds().Max.Y)) {
		return fmt.Errorf("Tile DrawGlobalXY error: global Y offset not within tile space (%d)", offsetY)
	}

	t.DrawLocalXY(src, offsetX, offsetY, config)

	return nil
}

// DrawLocation draws the provided image at the provided lat lng
func (t *Tile) DrawLocation(src image.Image, loc base.Location, config DrawConfig) {
	// Calculate location in pixel space
	x, y := MercatorLocationToPixel(loc.Latitude, loc.Longitude, t.Level, t.Size)
	t.DrawGlobalXY(src, int(x), int(y), config)
}
