package structs

// PublicRoom - a representation of the state of a room.
type PublicRoom struct {
	Building          string        `json:"-"`
	Room              string        `json:"-"`
	CurrentVideoInput string        `json:"currentVideoInput,omitempty"`
	CurrentAudioInput string        `json:"currentAudioInput,omitempty"`
	Power             string        `json:"power,omitempty"`
	Blanked           *bool         `json:"blanked,omitempty"`
	Muted             *bool         `json:"muted,omitempty"`
	Volume            *int          `json:"volume,omitempty"`
	Displays          []Display     `json:"displays,omitempty"`
	AudioDevices      []AudioDevice `json:"audioDevices,omitempty"`
}

// PublicDevice is a struct for inheriting
type PublicDevice struct {
	Name  string `json:"name,omitempty"`
	Power string `json:"power,omitempty"`
	Input string `json:"input,omitempty"`
}

// AudioDevice represents an audio device
type AudioDevice struct {
	PublicDevice
	Muted  *bool `json:"muted,omitempty"`
	Volume *int  `json:"volume,omitempty"`
}

// Display represents a display
type Display struct {
	PublicDevice
	Blanked *bool `json:"blanked,omitempty"`
}
