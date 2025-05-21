package domain

import "time"

type Calendar struct {
	CalendarId  int       `json:"calendar_id"`
	IsAvailable bool      `json:"is_available"`
	Date        time.Time `json:"date"`
	Hour        time.Time `json:"hour"`
}
