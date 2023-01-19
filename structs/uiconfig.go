package structs

// UIConfig - a representation of all the information needed to configure the touchpanel UI.
type UIConfig struct {
	ID                  string               `json:"_id,omitempty"`
	Api                 []string             `json:"api"`
	Panels              []Panel              `json:"panels"`
	Presets             []Preset             `json:"presets"`
	InputConfiguration  []IOConfiguration    `json:"inputConfiguration"`
	OutputConfiguration []IOConfiguration    `json:"outputConfiguration"`
	AudioConfiguration  []AudioConfiguration `json:"audioConfiguration"`
	PseudoInputs        []PseudoInput        `json:"pseudoInputs,omitempty"`
}

// Preset - a representation of what is controlled by this preset.
type Preset struct {
	Name                    string              `json:"name"`
	Icon                    string              `json:"icon"`
	Displays                []string            `json:"displays"`
	ShareablePresets        []string            `json:"shareablePresets"`
	ShareableDisplays       []string            `json:"shareableDisplays"`
	AudioDevices            []string            `json:"audioDevices"`
	Inputs                  []string            `json:"inputs"`
	IndependentAudioDevices []string            `json:"independentAudioDevices,omitempty"`
	AudioGroups             map[string][]string `json:"audioGroups,omitempty"`
	VolumeMatches           []string            `json:"volumeMatches,omitempty"`
	Commands                Commands            `json:"commands,omitempty"`
	Screens                 []string            `json:"screens"`
	Cameras                 []Camera            `json:"cameras"`
	Recording               Recording           `json:"recording"`
}

// Panel - a representation of a touchpanel and which preset it has.
type Panel struct {
	Hostname string   `json:"hostname"`
	UIPath   string   `json:"uipath"`
	Preset   string   `json:"preset"`
	Features []string `json:"features"`
}

// Commands - a representation of commands to be sent through the UI.
type Commands struct {
	PowerOn        []ConfigCommand `json:"powerOn,omitempty"`
	PowerOff       []ConfigCommand `json:"powerOff,omitempty"`
	InputSame      []ConfigCommand `json:"inputSame,omitempty"`
	InputDifferent []ConfigCommand `json:"inputDifferent,omitempty"`
	Delay          int             `json:"json:delay,omitempty"`
}

// ConfigCommand - ...I dunno, ask Danny.
type ConfigCommand struct {
	Method   string                 `json:"method"`
	Port     int                    `json:"port"`
	Endpoint string                 `json:"endpoint"`
	Body     map[string]interface{} `json:"body"`
}

// AudioConfiguration - a representation of how the audio is configured when using multiple displays.
type AudioConfiguration struct {
	Display      string   `json:"display"`
	AudioDevices []string `json:"audioDevices"`
	RoomWide     bool     `json:"roomWide"`
}

// IOConfiguration - a representation of an input or output device.
type IOConfiguration struct {
	Name        string            `json:"name"`
	Icon        string            `json:"icon"`
	Displayname *string           `json:"displayname,omitempty"`
	SubInputs   []IOConfiguration `json:"subInputs,omitempty"`
}

// PseudoInput - a fake input I guess
type PseudoInput struct {
	Displayname string `json:"displayname"`
	Config      []struct {
		Input   string   `json:"input"`
		Outputs []string `json:"outputs"`
	} `json:"config"`
}

// Template - the UI config and device list for a room for quick configuration.
type Template struct {
	ID          string   `json:"_id"`
	Description string   `json:"description"`
	UIConfig    UIConfig `json:"uiconfig"`
	BaseTypes   []string `json:"base_types"`
}

type Camera struct {
	DisplayName string `json:"displayName"`

	TiltUp      string `json:"tiltUp"`
	TiltDown    string `json:"tiltDown"`
	PanLeft     string `json:"panLeft"`
	PanRight    string `json:"panRight"`
	PanTiltStop string `json:"panTiltStop"`

	ZoomIn   string `json:"zoomIn"`
	ZoomOut  string `json:"zoomOut"`
	ZoomStop string `json:"zoomStop"`

	Presets []CameraPreset `json:"presets"`
}

type CameraPreset struct {
	DisplayName string `json:"displayName"`
	SetPreset   string `json:"setPreset"`
}

type Recording struct {
	Start   string `json:"start"`
	Stop    string `json:"stop"`
	MaxTime int    `json:"maxTime"`
}
