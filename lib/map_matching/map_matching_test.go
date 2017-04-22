/**
 * go-mapbox Map Matching Module
 * Wraps the mapbox Map Matching API for server side use
 * See https://www.mapbox.com/api-documentation/#retrieve-a-match for API information
 *
 * https://github.com/ryankurte/go-mapbox
 * Copyright 2017 Ryan Kurte
 */

package mapmatching

import (
	"os"
	"testing"
)

import (
	"github.com/ryankurte/go-mapbox/lib/base"
)

func TestMapMatching(t *testing.T) {

	token := os.Getenv("MAPBOX_TOKEN")
	if token == "" {
		t.Error("Mapbox API token not found")
		t.FailNow()
	}

	b := base.NewBase(token)

	MapMatching := NewMapMaptching(b)

	t.Run("Can Lookup Map Matching", func(t *testing.T) {
		timeStamps := []int64{1492878132, 1492878142, 1492878152, 1492878172, 1492878182, 1492878192, 1492878202, 1492878302}
		radiusList := []int{9, 6, 8, 11, 8, 4, 8, 8}
		var opts RequestOpts
		opts.SetGeometries(GeometryPolyline)
		opts.SetOverview(OverviewFull)
		opts.SetTimestamps(timeStamps)
		opts.SetSteps(false)
		opts.SetAnnotations([]AnnotationType{AnnotationDistance, AnnotationSpeed})
		opts.SetRadiuses(radiusList)

		locs := []base.Location{{37.75319556403746, -122.44254112243651}, {37.75373846204306, -122.44238018989562},
			{37.754111702111146, -122.44199395179749}, {37.75473941979767, -122.44177401065825},
			{37.755570713402115, -122.4412429332733}, {37.756401997666046, -122.44113564491273},
			{37.75677098309616, -122.44228899478911}, {37.756949113334784, -122.4424821138382}}

		res, err := MapMatching.GetMatching(locs, RoutingCycling, &opts)
		if err != nil {
			t.Error(err)
		}

		if Codes(res.Code) != CodeOK {
			t.Errorf("Invalid response code: %s", res.Code)
		}

	})
}
