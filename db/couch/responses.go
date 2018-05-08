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

type RoomQueryResponse struct {
	Docs     []room `json:"docs"`
	Bookmark string `json:"bookmark"`
	Warning  string `json:"warning"`
}

type roomConfiguration struct {
	structs.RoomConfiguration
	Rev string `json:"_rev,omitempty"`
}

type device struct {
	structs.Device
	Rev string `json:"_rev,omitempty"`
}

type deviceType struct {
	structs.DeviceType
	Rev string `json:"_rev,omitempty"`
}
