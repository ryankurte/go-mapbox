package mapbox

import (
	"os"
	"testing"
)

// Import the core module and any required APIs
import (
	"github.com/ryankurte/go-mapbox/lib/base"
	"github.com/ryankurte/go-mapbox/lib/directions"
	"github.com/ryankurte/go-mapbox/lib/directions_matrix"
	"github.com/ryankurte/go-mapbox/lib/geocode"
	"github.com/ryankurte/go-mapbox/lib/maps"
)

func TestMaps(t *testing.T) {
	// Fetch token from somewhere
	token := os.Getenv("MAPBOX_TOKEN")
	if token == "" {
		t.Errorf("No token found")
		t.FailNow()
	}

	// Create new mapbox instance
	mapBox := NewMapbox(token)

	// Map API
	_, _, err := mapBox.Maps.GetTile(maps.MapIDSatellite, 1, 0, 1, maps.MapFormatJpg90, true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Geocoding API

	// Forward Geocoding
	var forwardOpts geocode.ForwardRequestOpts
	forwardOpts.Limit = 1

	place := "2 lincoln memorial circle nw"

	_, err = mapBox.Geocode.Forward(place, &forwardOpts)
	if err != nil {
		t.Error(err)
	}

	// Reverse Geocoding
	var reverseOpts geocode.ReverseRequestOpts
	reverseOpts.Limit = 1

	loc := &base.Location{72.438939, 34.074122}

	_, err = mapBox.Geocode.Reverse(loc, &reverseOpts)
	if err != nil {
		t.Error(err)
	}

	// Directions API
	var directionOpts directions.RequestOpts

	locs := []base.Location{{-122.42, 37.78}, {-77.03, 38.91}}

	_, err = mapBox.Directions.GetDirections(locs, directions.RoutingCycling, &directionOpts)
	if err != nil {
		t.Error(err)
	}

	// Directions Matrix API
	var directionMatrixOpts directionsmatrix.RequestOpts
	// Only 1st and second points will act as a source the response will be a 2x3 matrix
	source := []string{"0", "1"}
	dest := []string{"all"}
	directionMatrixOpts.SetSources(source)
	directionMatrixOpts.SetDestinations(dest)

	points := []base.Location{{37.752759, -122.467600}, {37.762819, -122.460304}, {37.758095, -122.442253}}

	_, err = mapBox.DirectionsMatrix.GetDirectionsMatrix(points, directionsmatrix.RoutingCycling, &directionMatrixOpts)
	if err != nil {
		t.Error(err)
	}
}
