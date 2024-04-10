package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

type User struct {
	ID           string
	Name         string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

func (s *Store) SaveUser(user *User) error {
	query := `
		INSERT INTO users (id, name, access_token, refresh_token, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET name = $2, access_token = $3, refresh_token = $4, expires_at = $5
	`

	log.Printf("DBG s.db: %#v", s.db)
	log.Printf("DBG user: %#v", user)
	_, err := s.db.Exec(query, user.ID, user.Name, user.AccessToken, user.RefreshToken, user.ExpiresAt)
	return err
}

func (s *Store) GetUser(id string) (*User, error) {
	query := `
		SELECT id, name, access_token, refresh_token, expires_at
		FROM users
		WHERE id = $1
	`
	var user User
	err := s.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.AccessToken, &user.RefreshToken, &user.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

type Activity struct {
	ID            string    `db:"id"`
	UserID        string    `db:"user_id"`
	StartDate     time.Time `db:"start_date"`
	StartLat      float64   `db:"start_lat"`
	StartLng      float64   `db:"start_lng"`
	EndLat        float64   `db:"end_lat"`
	EndLng        float64   `db:"end_lng"`
	Distance      float64   `db:"distance"`
	MovingTime    int       `db:"moving_time"`
	ElapsedTime   int       `db:"elapsed_time"`
	ElevLow       float64   `db:"elev_low"`
	ElevHigh      float64   `db:"elev_high"`
	TotalElevGain float64   `db:"total_elevation_gain"`
}

func (s *Store) SaveActivity(activity *Activity) error {
	_, err := s.db.Exec(`
        INSERT INTO activities (
            id, user_id, start_date, start_lat, start_lng, end_lat, end_lng, distance, moving_time,
            elapsed_time, elev_low, elev_high, total_elevation_gain
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `, activity.ID, activity.UserID, activity.StartDate, activity.StartLat, activity.StartLng, activity.EndLat, activity.EndLng, activity.Distance,
		activity.MovingTime, activity.ElapsedTime, activity.ElevLow, activity.ElevHigh, activity.TotalElevGain)

	return err
}

func (s *Store) GetActivity(id string) (*Activity, error) {
	var activity Activity

	err := s.db.QueryRow(`
        SELECT id, user_id, start_lat, start_lng, end_lat, end_lng, end_latlng, distance, moving_time,
               elapsed_time, elev_low, elev_high, total_elevation_gain
        FROM activities
        WHERE id = $1
    `, id).Scan(&activity.ID, &activity.UserID, &activity.StartDate, &activity.StartLat, &activity.StartLng, &activity.EndLat, &activity.EndLng, &activity.Distance,
		&activity.MovingTime, &activity.ElapsedTime, &activity.ElevLow, &activity.ElevHigh, &activity.TotalElevGain)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &activity, nil
}

func (s *Store) ListActivities(userID string) ([]*Activity, error) {
	rows, err := s.db.Query(`
        SELECT id, user_id, start_date, start_lat, start_lng, end_lat, end_lng, distance, moving_time,
               elapsed_time, elev_low, elev_high, total_elevation_gain
        FROM activities
        WHERE user_id = $1
        ORDER BY start_date DESC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*Activity
	for rows.Next() {
		var activity Activity
		err := rows.Scan(&activity.ID, &activity.UserID, &activity.StartDate, &activity.StartLat, &activity.StartLng, &activity.EndLat, &activity.EndLng, &activity.Distance,
			&activity.MovingTime, &activity.ElapsedTime, &activity.ElevLow, &activity.ElevHigh, &activity.TotalElevGain)
		if err != nil {
			return nil, err
		}
		activities = append(activities, &activity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}

func ConvertStravaActivityToActivity(stravaActivity *StravaActivity) *Activity {
	activity := &Activity{
		ID:            fmt.Sprintf("%d", stravaActivity.ID),
		UserID:        fmt.Sprintf("%d", stravaActivity.Athlete.ID),
		StartDate:     stravaActivity.StartDate,
		StartLat:      stravaActivity.StartLanLng[0],
		StartLng:      stravaActivity.StartLanLng[1],
		EndLat:        stravaActivity.EndLanLng[0],
		EndLng:        stravaActivity.EndLanLng[1],
		Distance:      stravaActivity.Distance,
		MovingTime:    stravaActivity.MovingTime,
		ElapsedTime:   stravaActivity.ElapsedTime,
		ElevLow:       stravaActivity.ElevLow,
		ElevHigh:      stravaActivity.ElevHigh,
		TotalElevGain: stravaActivity.TotalElevGain,
	}

	return activity
}
