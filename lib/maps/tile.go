package maps

import (
	"fmt"
	"image"
	"image/draw"
	"log"
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
	log.Printf("Pixel: %+v", p)
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

func GradientInterpolate2D(m [][]float64, ignore float64) [][]float64 {
	lenY := len(m)
	lenX := len(m[0])

	log.Printf("Original:      %+v", m)

	interpolatedX := make([][]float64, lenY)
	for i := range m {
		interpolatedX[i] = make([]float64, lenX)
	}

	for y := 0; y < lenY; y++ {
		interpolateCount := 0
		lastX := 0

		for x := 0; x < lenX; x++ {
			p := m[y][x]
			log.Printf("X: %d, y: %d, p: %.2f ic: %d", x, y, p, interpolateCount)
			if p != ignore {
				interpolatedX[y][x] = p

				if interpolateCount > 0 && !(lastX == 0 && m[y][0] == ignore) {
					last, next := m[y][lastX], p
					delta := (next - last) / float64(interpolateCount+1)
					log.Printf("Interpolate H last X: %d (%.2f) next X: %d (%.2f) count: %d delta: %.2f",
						lastX, last, x, next, interpolateCount, delta)

					for j := 1; j < interpolateCount+1; j++ {
						val := last + float64(j)*delta
						log.Printf("Filling x: %d value: %.2f", lastX+j, val)
						interpolatedX[y][lastX+j] = val
					}

					interpolateCount = 0
				}
				lastX = x

			} else {
				if !(lastX == 0 && m[y][0] == ignore) {
					interpolateCount++
				}

			}
		}
	}

	log.Printf("Interpolated X:  %+v", interpolatedX)

	interpolatedY := make([][]float64, lenY)
	for i := range m {
		interpolatedY[i] = make([]float64, lenX)
	}

	for x := 0; x < lenX; x++ {
		interpolateCount := 0
		lastY := 0

		for y := 0; y < lenY; y++ {
			p := m[y][x]
			log.Printf("Vertical iterator X: %d, y: %d, p: %.2f ic: %d", x, y, p, interpolateCount)
			if p != ignore {
				interpolatedY[y][x] = p

				if interpolateCount > 0 && !(lastY == 0 && m[0][x] == ignore) {
					last, next := m[lastY][x], p
					delta := (next - last) / float64(interpolateCount+1)
					log.Printf("Interpolate V last Y: %d (%.2f) next Y: %d (%.2f) count: %d delta: %.2f",
						lastY, last, y, next, interpolateCount, delta)

					for j := 1; j < interpolateCount+1; j++ {
						val := last + float64(j)*delta
						log.Printf("Filling y: %d value: %.2f", lastY+j, val)
						interpolatedY[lastY+j][x] = val
					}

					interpolateCount = 0
				}
				lastY = y

			} else {
				if !(lastY == 0 && m[0][x] == ignore) {
					interpolateCount++
				}

			}
		}
	}

	log.Printf("Interpolated Y:  %+v", interpolatedY)

	interpolated := make([][]float64, lenY)
	for y := 0; y < lenY; y++ {
		interpolated[y] = make([]float64, lenX)
		for x := 0; x < lenX; x++ {
			intX := interpolatedX[y][x]
			intY := interpolatedY[y][x]
			if intX == ignore {
				interpolated[y][x] = intY
			} else if intY == ignore {
				interpolated[y][x] = intX
			} else {
				interpolated[y][x] = (interpolatedX[y][x] + interpolatedY[y][x]) / 2
			}
		}
	}

	log.Printf("Interpolated 2D: %+v", interpolated)

	return interpolated
}
