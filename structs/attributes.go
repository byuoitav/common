package structs

// MenuTree is a wrapper for the list of groups
type MenuTree struct {
	Groups []Group `json:"groups"`
}

// Group is a collection of attribute presets to create devices that fall into this group
type Group struct {
	ID        string         `json:"_id"`
	Icon      string         `json:"icon,omitempty"`
	Subgroups []Group        `json:"sub-groups,omitempty"`
	Presets   []AttributeSet `json:"presets,omitempty"`
}

// AttributeSet is an object that contains a set of attributes and an identifier for this set
type AttributeSet struct {
	Name       string                 `json:"name"`
	DeviceType string                 `json:"device-type"`
	DeviceName string                 `json:"device-name,omitempty"`
	DeviceIcon string                 `json:"device-icon,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
}
