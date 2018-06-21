package structs

type UIConfig struct {
	ID                  string               `json:"_id,omitempty"`
	Api                 []string             `json:"api"`
	Panels              []Panel              `json:"panels"`
	Presets             []Preset             `json:"presets"`
	InputConfiguration  []IOConfiguration    `json:"inputConfiguration"`
	OutputConfiguration []IOConfiguration    `json:"outputConfiguration"`
	AudioConfiguration  []AudioConfiguration `json:"audioConfiguration"`
}

type Preset struct {
	Name                    string   `json:"name"`
	Icon                    string   `json:"icon"`
	Displays                []string `json:"displays"`
	ShareableDisplays       []string `json:"shareableDisplays"`
	AudioDevices            []string `json:"audioDevices"`
	Inputs                  []string `json:"inputs"`
	IndependentAudioDevices []string `json:"independentAudioDevices"`
}

type Panel struct {
	Hostname string   `json:"hostname"`
	UIPath   string   `json:"uipath"`
	Preset   string   `json:"preset"`
	Features []string `json:"features"`
}

type AudioConfiguration struct {
	Display      string   `json:"display"`
	AudioDevices []string `json:"audioDevices"`
	RoomWide     bool     `json:"roomWide"`
}

type IOConfiguration struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}
