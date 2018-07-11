package structs

// PowerStatus is a base evaluator struct.
type PowerStatus struct {
	Power string `json:"power"`
}

// BlankedStatus is a base evaluator struct.
type BlankedStatus struct {
	Blanked bool `json:"blanked"`
}

// MuteStatus is a base evaluator struct.
type MuteStatus struct {
	Muted bool `json:"muted"`
}

// InputStatus is a base evaluator struct.
type InputStatus struct {
	Input string `json:"input,omitempty"`
}

// VolumeStatus is a base evaluator struct.
type VolumeStatus struct {
	Volume int `json:"volume"`
}
