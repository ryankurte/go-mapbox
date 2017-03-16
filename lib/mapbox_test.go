package mapbox

import (
	"os"
	"testing"
)

func MapboxTest(t *testing.T) {

	token := os.Getenv("MAPBOX_TOKEN")
	if token == "" {
		t.Error("Mapbox API token not found")
		t.FailNow()
	}

	t.Run("Can make API requests", func(t *testing.T) {

	})

	t.Run("Can reverse geocode", func(t *testing.T) {

	})

}
