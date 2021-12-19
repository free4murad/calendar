package models

import "time"

// Calendar represents a calendar of a person along with a bunch of meetings
type Calendar struct {
	Name     string    `json:"name"`
	Meetings []Meeting `json:"meetings"`
}

// Slot represents a time slot with start and end time
type Slot struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// Meeting represents a slot with a subject
type Meeting struct {
	Slot
	Subject string `json:"subject"`
}
