package structs

import "regexp"

// Jobs .
type Jobs struct {
	Jobs []JobConfig `json:"jobs"`
}

// JobConfig defines a configuration of a specific job.
type JobConfig struct {
	Name     string      `json:"name"`
	Triggers []Trigger   `json:"triggers"`
	Enabled  bool        `json:"enabled"`
	Context  interface{} `json:"context"`
}

// Trigger matches something that causes a job to be ran.
type Trigger struct {
	Type  string       `json:"type"`            // required for all
	At    *string      `json:"at,omitempty"`    // required for 'time'
	Every *string      `json:"every,omitempty"` // required for 'interval'
	Match *MatchConfig `json:"match,omitempty"` // required for 'event'
}

// MatchConfig contains the logic for building/matching regex for events that come in
type MatchConfig struct {
	Count int

	GeneratingSystem string   `json:"generating-system"`
	Timestamp        string   `json:"timestamp"`
	EventTags        []string `json:"event-tags"`
	Key              string   `json:"key"`
	Value            string   `json:"value"`
	User             string   `json:"user"`
	Data             string   `json:"data,omitempty"`
	AffectedRoom     struct {
		BuildingID string `json:"buildingID,omitempty"`
		RoomID     string `json:"roomID,omitempty"`
	} `json:"affected-room"`
	TargetDevice struct {
		BuildingID string `json:"buildingID,omitempty"`
		RoomID     string `json:"roomID,omitempty"`
		DeviceID   string `json:"deviceID,omitempty"`
	} `json:"target-device"`

	Regex struct {
		GeneratingSystem *regexp.Regexp
		Timestamp        *regexp.Regexp
		EventTags        []*regexp.Regexp
		Key              *regexp.Regexp
		Value            *regexp.Regexp
		User             *regexp.Regexp
		Data             *regexp.Regexp
		AffectedRoom     struct {
			BuildingID *regexp.Regexp
			RoomID     *regexp.Regexp
		}
		TargetDevice struct {
			BuildingID *regexp.Regexp
			RoomID     *regexp.Regexp
			DeviceID   *regexp.Regexp
		}
	}
}
