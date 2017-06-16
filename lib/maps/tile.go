package maps

import (
	"fmt"
	"image"
	"image/draw"
	"math"

	"github.com/ryankurte/go-mapbox/lib/base"
	"image/color"
)

// Tile is a wrapper around an image that includes positioning data
type Tile struct {
	image.Image
	Level uint64 // Tile zoom level
	Size  uint64 // Tile size
	X, Y  uint64 // Tile X and Y postions (Web Mercurator projection)
}

const (
	SizeStandard uint64 = 256
	SizeHighDPI  uint64 = 512
)

// NewTile creates a tile with a base RGBA object
func NewTile(x, y, level, size uint64, src image.Image) Tile {
	// Convert image to RGBA
	b := src.Bounds()
	m := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), src, b.Min, draw.Src)

	return Tile{
		Image: m,
		X:     x,
		Y:     y,
		Level: level,
		Size:  size,
	}
}

// LocationToPixel translates a global location to a pixel on the tile
func (t *Tile) LocationToPixel(loc base.Location) (float64, float64, error) {
	x, y := MercatorLocationToPixel(loc.Latitude, loc.Longitude, t.Level, t.Size)
	offsetX, offsetY := x-float64(t.X*t.Size), y-float64(t.Y*t.Size)

	if xMax := float64(t.Image.Bounds().Max.X); (offsetX < 0) || (offsetX > xMax) {
		return 0, 0, fmt.Errorf("Tile LocationToPixel error: global X offset not within tile space (x: %d max: %d)", offsetX, int(xMax))
	}
	if yMax := float64(t.Image.Bounds().Max.Y); (offsetY < 0) || (offsetY > yMax) {
		return 0, 0, fmt.Errorf("Tile LocationToPixel error: global Y offset not within tile space (y: %d max: %d)", offsetY, int(yMax))
	}

	return offsetX, offsetY, nil
}

// PixelToLocation translates a pixel location in the tile into a global location
func (t *Tile) PixelToLocation(x, y float64) (*base.Location, error) {
	if xMax := float64(t.Image.Bounds().Max.X); (x < 0) || (x > xMax) {
		return nil, fmt.Errorf("Tile LocationToPixel error: global X offset not within tile space (x: %.2f max: %d)", x, int(xMax))
	}
	if yMax := float64(t.Image.Bounds().Max.Y); (y < 0) || (y > yMax) {
		return nil, fmt.Errorf("Tile LocationToPixel error: global Y offset not within tile space (y: %.2f max: %d)", y, int(yMax))
	}

	offsetX, offsetY := x+float64(t.X*t.Size), y+float64(t.Y*t.Size)
	lat, lng := MercatorPixelToLocation(offsetX, offsetY, t.Level, t.Size)

	return &base.Location{Latitude: lat, Longitude: lng}, nil
}

// Justify sets image offsets for drawing
type Justify string

// Constant justification types
const (
	JustifyTop    Justify = "top"
	JustifyLeft   Justify = "left"
	JustifyCenter Justify = "center"
	JustifyBottom Justify = "bottom"
	JustifyRight  Justify = "right"
)

// DrawConfig configures image drawing
type DrawConfig struct {
	Vertical   Justify
	Horizontal Justify
}

// Center preconfigured centering helper
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

func (t *Tile) translateGlobalToLocalXY(x, y int) (int, int, error) {
	offsetX := x - int(t.X*t.Size)
	offsetY := y - int(t.Y*t.Size)

	if xMax := int(t.Image.Bounds().Max.X); (offsetX < 0) || (offsetX > xMax) {
		return 0, 0, fmt.Errorf("Tile DrawGlobalXY error: global X offset not within tile space (x: %d max: %d)", offsetX, xMax)
	}
	if yMax := int(t.Image.Bounds().Max.Y); (offsetY < 0) || (offsetY > yMax) {
		return 0, 0, fmt.Errorf("Tile DrawGlobalXY error: global Y offset not within tile space (y: %d max: %d)", offsetY, yMax)
	}

	return offsetX, offsetY, nil
}

// DrawGlobalXY draws the provided image at the global X/Y coordinates
func (t *Tile) DrawGlobalXY(src image.Image, x, y int, config DrawConfig) error {
	offsetX, offsetY, err := t.translateGlobalToLocalXY(x, y)
	if err != nil {
		return err
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

// Interpolate function to be passed to generic line interpolator
type Interpolate func(pixel color.Color) color.Color

// InterpolateLocalXY interpolates a line between two local points and calls the interpolate function on each point
func (t *Tile) InterpolateLocalXY(x1, y1, x2, y2 int, interpolate Interpolate) {
	// This is a bit insane because the ordering between (x1, y1) and (x2, y2) must be preserved
	// So that points

	dx := int(float64(x2) - float64(x1))
	dy := int(float64(y2) - float64(y1))

	len := int(math.Sqrt(math.Pow(float64(dx), 2) + math.Pow(float64(dy), 2)))

	img := t.Image.(draw.Image)

	for i := 0; i < len; i++ {
		x := x1 + i*dx/len
		y := y1 + i*dy/len
		pixel := interpolate(img.At(x, y))
		img.Set(x, y, pixel)
	}
}

// InterpolateGlobalXY interpolates a line between two global points and calls the interpolate function on each point
func (t *Tile) InterpolateGlobalXY(x1, y1, x2, y2 int, interpolate Interpolate) error {
	offsetX1, offsetY1, err := t.translateGlobalToLocalXY(x1, y1)
	if err != nil {
		return err
	}
	offsetX2, offsetY2, err := t.translateGlobalToLocalXY(x2, y2)
	if err != nil {
		return err
	}
	t.InterpolateLocalXY(offsetX1, offsetY1, offsetX2, offsetY2, interpolate)
	return nil
}

// InterpolateLocations interpolates a line between two locations and calls the interpolate function on each point
func (t *Tile) InterpolateLocations(loc1, loc2 base.Location, interpolate Interpolate) error {
	x1, y1 := MercatorLocationToPixel(loc1.Latitude, loc1.Longitude, t.Level, t.Size)
	x2, y2 := MercatorLocationToPixel(loc2.Latitude, loc2.Longitude, t.Level, t.Size)
	return t.InterpolateGlobalXY(int(x1), int(y1), int(x2), int(y2), interpolate)
}

// DrawLine uses InterpolateLocations to draw a line between two points
func (t *Tile) DrawLine(loc1, loc2 base.Location, c color.Color) {
	t.InterpolateLocations(loc1, loc2, func(color.Color) color.Color {
		return c
	})
}

func (t *Tile) GetAltitude(loc base.Location) (float64, error) {
	x, y := MercatorLocationToPixel(loc.Latitude, loc.Longitude, t.Level, t.Size)
	offsetX, offsetY, err := t.translateGlobalToLocalXY(int(x), int(y))
	if err != nil {
		return 0.0, err
	}
	p := t.Image.At(offsetX, offsetY).(color.RGBA)
	return PixelToHeight(p.R, p.G, p.B), nil
}

func (t *Tile) InterpolateAltitudes(loc1, loc2 base.Location) []float64 {
	altitudes := make([]float64, 0)
	t.InterpolateLocations(loc1, loc2, func(c color.Color) color.Color {
		p := c.(color.RGBA)
		alt := PixelToHeight(p.R, p.G, p.B)
		altitudes = append(altitudes, alt)
		return c
	})
	return altitudes
}

func (t *Tile) GetHighestAltitude() float64 {
	p := t.Image.At(0, 0).(color.RGBA)
	max := PixelToHeight(p.R, p.G, p.B)
	for y := 0; y < t.Image.Bounds().Dy(); y++ {
		for x := 0; x < t.Image.Bounds().Dx(); x++ {
			p := t.Image.At(x, y).(color.RGBA)
			alt := PixelToHeight(p.R, p.G, p.B)
			if alt > max {
				max = alt
			}
		}
	}
	return max
}

func (t *Tile) FlattenAltitudes(maxHeight float64) Tile {
	img := image.NewRGBA(t.Image.Bounds())

	for y := 0; y < t.Image.Bounds().Dy(); y++ {
		for x := 0; x < t.Image.Bounds().Dx(); x++ {
			p := t.Image.At(x, y).(color.RGBA)
			alt := uint8(PixelToHeight(p.R, p.G, p.B) / maxHeight * 255)
			img.Set(x, y, color.RGBA{R: alt, G: alt, B: alt, A: 255})
		}
	}

	return NewTile(t.X, t.Y, t.Level, t.Size, img)
}
