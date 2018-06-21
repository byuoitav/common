package couch

import "github.com/byuoitav/common/structs"

type building struct {
	Rev string `json:"_rev,omitempty"`
	*structs.Building
}

type buildingQueryResponse struct {
	Docs     []building `json:"docs"`
	Bookmark string     `json:"bookmark"`
	Warning  string     `json:"warning"`
}

type room struct {
	Rev string `json:"_rev,omitempty"`
	*structs.Room
}

type roomQueryResponse struct {
	Docs     []room `json:"docs"`
	Bookmark string `json:"bookmark"`
	Warning  string `json:"warning"`
}

type roomConfiguration struct {
	Rev string `json:"_rev,omitempty"`
	*structs.RoomConfiguration
}

type roomConfigurationQueryResponse struct {
	Docs     []roomConfiguration `json:"docs"`
	Bookmark string              `json:"bookmark"`
	Warning  string              `json:"warning"`
}

type device struct {
	Rev string `json:"_rev,omitempty"`
	*structs.Device
}

type deviceQueryResponse struct {
	Docs     []device `json:"docs"`
	Bookmark string   `json:"bookmark"`
	Warning  string   `json:"warning"`
}

type deviceType struct {
	Rev string `json:"_rev,omitempty"`
	*structs.DeviceType
}

type deviceTypeQueryResponse struct {
	Docs     []deviceType `json:"docs"`
	Bookmark string       `json:"bookmark"`
	Warning  string       `json:"warning"`
}

type uiconfig struct {
	Rev string `json:"_rev",omitempty`
	*structs.UIConfig
}

type icons struct {
	Rev      string   `json:"_rev",omitempty`
	IconList []string `json:"Icons"`
}

type deviceRoles struct {
	Rev      string   `json:"_rev",omitempty`
	RoleList []string `json:"roles"`
}

type roomDesignations struct {
	Rev       string   `json:"_rev",omitempty`
	DesigList []string `json:"designations"`
}

type tags struct {
	Rev     string   `json:"_rev",omitempty`
	TagList []string `json:"tags"`
}
