package couch

import "github.com/byuoitav/common/structs"

type building struct {
	structs.Building
	Rev string `json:"_rev,omitempty"`
}

type buildingQueryResponse struct {
	Docs     []building `json:"docs"`
	Bookmark string     `json:"bookmark"`
	Warning  string     `json:"warning"`
}

type room struct {
	structs.Room
	Rev string `json:"_rev,omitempty"`
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
