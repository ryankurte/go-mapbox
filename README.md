# go-mapbox

Mapbox API wrappers for Golang

[![Documentation](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/ryankurte/go-mapbox/lib)
[![GitHub tag](https://img.shields.io/github/tag/ryankurte/go-mapbox.svg)](https://github.com/ryankurte/go-mapbox)
[![Build Status](https://travis-ci.org/ryankurte/go-mapbox.svg?branch=master)](https://travis-ci.org/ryankurte/go-mapbox)

See [here](https://golanglibs.com/top?q=mapbox) for other golang/mapbox projects.

## Status

Very early WIP, pull requests and issues are most welcome. See [lib/geocode/](lib/geocode) or [lib/directions/](lib/directions) for an example module to mimic.

### Modules

- [X] Geocoding
- [X] Directions
- [ ] Styles
- [X] Maps
- [ ] Static
- [ ] Datasets

## Examples

### Initialisation

```go
// Import the core module (and any required APIs)
import (
    "gopkg.in/ryankurte/go-mapbox.v0/lib"
    "gopkg.in/ryankurte/go-mapbox.v0/lib/base"
)

// Fetch token from somewhere
token := os.Getenv("MAPBOX_TOKEN")

// Create new mapbox instance
mapBox := mapbox.NewMapbox(token)

```

### Map API
``` go
import (
    "gopkg.in/ryankurte/go-mapbox.v0/lib/maps"
)

img, err := mapBox.Maps.GetTiles(maps.MapIDSatellite, 1, 0, 0, maps.MapFormatJpg90, true)
```

### Geocoding

```go
import (
    "gopkg.in/ryankurte/go-mapbox.v0/lib/geocode"
)

// Forward Geocoding
var forwardOpts geocode.ForwardRequestOpts
forwardOpts.Limit = 1

place := "2 lincoln memorial circle nw"

forward, err := mapBox.Geocode.Forward(place, &forwardOpts)


// Reverse Geocoding
var reverseOpts geocode.ReverseRequestOpts
reverseOpts.Limit = 1

loc := &base.Location{72.438939, 34.074122}

reverse, err := mapBox.Geocode.Reverse(loc, &reverseOpts)

```

### Directions

```go
import (
    "gopkg.in/ryankurte/go-mapbox.v0/lib/directions"
)

var directionOpts directions.RequestOpts

locs := []base.Location{{-122.42, 37.78}, {-77.03, 38.91}}

directions, err := mapBox.Directions.GetDirections(locs, directions.RoutingCycling, &directionOpts)

```

## Layout

- [lib/base](lib/base/) contains a common base for API modules
- [lib/maps](lib/maps/) contains the maps API module
- [lib/directions](lib/directions/) contains the directions API module
- [lib/geocode](lib/geocode/) contains the geocoding API module

---

If you have any questions, comments, or suggestions, feel free to open an issue or a pull request.

