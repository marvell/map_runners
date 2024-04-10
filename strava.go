package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Strava struct {
	oauth2Config *oauth2.Config
}

func (app *Application) InitStrava() {
	app.strava = &Strava{
		oauth2Config: &oauth2.Config{
			ClientID:     app.cfg.StravaClientID,
			ClientSecret: app.cfg.StravaClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.strava.com/oauth/authorize",
				TokenURL: "https://www.strava.com/oauth/token",
			},
			RedirectURL: app.cfg.StravaRedirectURL,
			Scopes:      []string{"activity:read_all"},
		},
	}
}

func (app *Application) StravaLoginHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		url := app.strava.oauth2Config.AuthCodeURL("state")
		c.Redirect(http.StatusFound, url)
	}
}

func (app *Application) StravaCallbackHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		code := c.Query("code")
		token, err := app.strava.oauth2Config.Exchange(c.Request.Context(), code)
		if err != nil {
			printError(c, err)
			return
		}
		log.Printf("DBG token: %+v\n", token)

		athlete := token.Extra("athlete").(map[string]interface{})
		log.Printf("DBG athlete: %+v\n", athlete)
		user := &User{
			ID:           fmt.Sprintf("%.0f", athlete["id"]),
			Name:         fmt.Sprintf("%s %s", athlete["firstname"], athlete["lastname"]),
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    token.Expiry,
		}
		log.Printf("DBG user: %+v\n", user)

		if err := app.store.SaveUser(user); err != nil {
			printError(c, err)
			return
		}

		cookie := &http.Cookie{
			Name:     "user_id",
			Value:    user.ID,
			Path:     "/",
			Domain:   "localhost:8000",
			Expires:  time.Now().Add(time.Minute),
			Secure:   true,
			HttpOnly: true,
		}
		c.SetCookie(cookie.Name, cookie.Value, int(cookie.Expires.Unix()), cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)

		go func() {
			err := app.syncActivities(user.ID)
			if err != nil {
				log.Printf("ERR sync activities: %s", err)
			}
		}()

		c.Redirect(http.StatusFound, fmt.Sprintf("/%s", user.ID))
	}
}

type StravaActivity struct {
	ID      int64 `json:"id"`
	Athlete struct {
		ID int64 `json:"id"`
	} `json:"athlete"`

	SportType StravaSportType `json:"sport_type"`

	StartDate time.Time `json:"start_date"`

	StartLanLng [2]float64 `json:"start_latlng"`
	EndLanLng   [2]float64 `json:"end_latlng"`

	Distance    float64 `json:"distance"`
	MovingTime  int     `json:"moving_time"`
	ElapsedTime int     `json:"elapsed_time"`

	ElevLow       float64 `json:"elev_low"`
	ElevHigh      float64 `json:"elev_high"`
	TotalElevGain float64 `json:"total_elevation_gain"`
}

type StravaSportType string

const (
	RunStravaSportType StravaSportType = "Run"
)

func getActivities(accessToken string, from, to time.Time) ([]*StravaActivity, error) {
	var activities []*StravaActivity

	for i := 1; ; i++ {
		params := url.Values{}
		params.Add("access_token", accessToken)
		params.Add("per_page", "100")
		params.Add("page", fmt.Sprint(i))
		params.Add("after", fmt.Sprint(from.Unix()))
		params.Add("before", fmt.Sprint(to.Unix()))

		req, err := http.NewRequest("GET", "https://www.strava.com/api/v3/athlete/activities?"+params.Encode(), nil)
		if err != nil {
			return nil, err
		}

		reqBody, err := httputil.DumpRequestOut(req, false)
		if err != nil {
			return nil, err
		}
		log.Printf("DBG request: %s", reqBody)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		respBody, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, err
		}
		log.Printf("DBG response: %s", respBody)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var v []*StravaActivity
		err = json.Unmarshal(body, &v)
		if err != nil {
			return nil, err
		}

		if len(v) == 0 {
			break
		}

		activities = append(activities, v...)

		time.Sleep(5 * time.Second)
	}

	return activities, nil
}

func (app *Application) syncActivities(userID string) error {
	user, err := app.store.GetUser(userID)
	if err != nil {
		return err
	}

	now := time.Now()

	activities, err := getActivities(user.AccessToken, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), now)
	if err != nil {
		return err
	}

	for _, a := range activities {
		if a.SportType != RunStravaSportType {
			continue
		}

		if a.StartLanLng[0] == 0 || a.StartLanLng[1] == 0 {
			continue
		}

		err = app.store.SaveActivity(ConvertStravaActivityToActivity(a))
		if err != nil {
			return err
		}
	}

	return nil
}
