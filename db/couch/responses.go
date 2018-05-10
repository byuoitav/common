package couch

import "github.com/byuoitav/common/structs"

type building struct {
	Rev string `json:"_rev,omitempty"`
	structs.Building
}

type buildingQueryResponse struct {
	Docs     []building `json:"docs"`
	Bookmark string     `json:"bookmark"`
	Warning  string     `json:"warning"`
}

type room struct {
	devices []device `json:"devices"`
	Rev     string   `json:"_rev,omitempty"`
	structs.Room
}

type roomQueryResponse struct {
	Docs     []room `json:"docs"`
	Bookmark string `json:"bookmark"`
	Warning  string `json:"warning"`
}

type roomConfiguration struct {
	structs.RoomConfiguration
	Rev string `json:"_rev,omitempty"`
}

type roomConfigurationQueryResponse struct {
	Docs     []roomConfiguration `json:"docs"`
	Bookmark string              `json:"bookmark"`
	Warning  string              `json:"warning"`
}

type device struct {
	structs.Device
	Rev string `json:"_rev,omitempty"`
}

type deviceQueryResponse struct {
	Docs     []device `json:"docs"`
	Bookmark string   `json:"bookmark"`
	Warning  string   `json:"warning"`
}

type deviceType struct {
	structs.DeviceType
	Rev string `json:"_rev,omitempty"`
}

type deviceTypeQueryResponse struct {
	Docs     []deviceType `json:"docs"`
	Bookmark string       `json:"bookmark"`
	Warning  string       `json:"warning"`
}
