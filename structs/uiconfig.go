package structs

type UIConfig struct {
	ID                  string               `json:"_id,omitempty"`
	Api                 []string             `json:"api"`
	Panels              []Panel              `json:"panels"`
	Presets             []Preset             `json:"presets"`
	InputConfiguration  []IOConfiguration    `json:"inputConfiguration"`
	OutputConfiguration []IOConfiguration    `json:"outputConfiguration,omitempty"`
	AudioConfiguration  []AudioConfiguration `json:"audioConfiguration"`
	PseudoInputs        []PseudoInput        `json:"pseudoInputs"`
}

type Preset struct {
	Name                    string   `json:"name"`
	Icon                    string   `json:"icon"`
	Displays                []string `json:"displays"`
	ShareableDisplays       []string `json:"shareableDisplays,omitempty"`
	AudioDevices            []string `json:"audioDevices"`
	Inputs                  []string `json:"inputs"`
	IndependentAudioDevices []string `json:"independentAudioDevices,omitempty"`
	Commands                struct {
		PowerOn  []ConfigCommand `json:"powerOn"`
		PowerOff []ConfigCommand `json:"powerOff"`
	} `json:"commands"`
}

type Panel struct {
	Hostname string   `json:"hostname"`
	UIPath   string   `json:"uipath"`
	Preset   string   `json:"preset"`
	Features []string `json:"features"`
}

type ConfigCommand struct {
	Method   string                 `json:"method"`
	Port     int                    `json:"port"`
	Endpoint string                 `json:"endpoint"`
	Body     map[string]interface{} `json:"body"`
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

type PseudoInput struct {
	Displayname string `json:"displayname"`
	Config      []struct {
		Input   string   `json:"input"`
		Outputs []string `json:"outputs"`
	} `json:"config"`
}
