package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

func (app *Application) IndexHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		userID, err := c.Cookie("user_id")
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			printError(c, err)
			return
		}

		if userID != "" {
			c.Redirect(http.StatusFound, fmt.Sprintf("/%s", userID))
			return
		}

		c.Redirect(http.StatusFound, "/strava/login")
	}
}

type CountryActivitiesSummary struct {
	ActivityNumber int
	SumDistance    float64
}

func (app *Application) UserHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		userID := c.Param("user_id")

		activities, err := app.store.ListActivities(userID)
		if err != nil {
			printError(c, err)
			return
		}

		countriesGeoJson, err := os.ReadFile("./static/ne_50m_admin_0_countries.json")
		if err != nil {
			printError(c, err)
			return
		}

		featureCollection, err := geojson.UnmarshalFeatureCollection(countriesGeoJson)
		if err != nil {
			printError(c, err)
			return
		}

		countries := map[string]*CountryActivitiesSummary{}

		for _, a := range activities {
			startPoint := orb.Point{a.StartLng, a.StartLat}

			for _, feature := range featureCollection.Features {
				contains := false
				switch v := feature.Geometry.(type) {
				case orb.Polygon:
					if planar.PolygonContains(v, startPoint) {
						contains = true
					}
				case orb.MultiPolygon:
					if planar.MultiPolygonContains(v, startPoint) {
						contains = true
					}
				}

				if contains {
					countryName := feature.Properties.MustString("ADM0_A3")

					if _, found := countries[countryName]; !found {
						countries[countryName] = &CountryActivitiesSummary{}
					}

					countries[countryName].ActivityNumber++
					countries[countryName].SumDistance += a.Distance
				}
			}
		}

		c.JSON(200, gin.H{
			"countries": countries,
		})
	}
}
