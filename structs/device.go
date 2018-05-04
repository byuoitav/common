package structs

import "strings"

type Device struct {
	ID          string     `json:"_id"`
	Rev         string     `json:"_rev,omitempty"`
	Address     string     `json:"address"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DisplayName string     `json:"display-name"`
	Type        DeviceType `json:"type"`
	Class       string     `json:"class"`
	Roles       []Role     `json:"roles"`
	Ports       []Port     `json:"ports"`
}

type DeviceType struct {
	ID          string       `json:"_id"`
	Rev         string       `json:"_rev,omitempty"`
	Name        string       `json:"name"`
	Class       string       `json:"class"`
	Description string       `json:"description"`
	Ports       []Port       `json:"ports"`
	PowerStates []PowerState `json:"power-states"`
	Commands    []Command    `json:"commands"`
}

type PowerState struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Port struct {
	ID                string `json:"_id"`
	FriendlyName      string `json:"friendly-name"`
	Name              string `json:"name"`
	SourceDevice      string `json:"source-device"`
	DestinationDevice string `json:"destination-device"`
	Description       string `json:"description"`
	PortType          string `json:"port-type"`
}

type Role struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Command struct {
	BaseInfo
	Microservice Microservice `json:"microservice"`
	Endpoint     Endpoint     `json:"endpoint"`
}

type DeviceQueryResponse struct {
	Docs     []Device `json:"docs"`
	Bookmark string   `json:"bookmark"`
	Warning  string   `json:"warning"`
}

type Microservice struct {
	BaseInfo
	Address string `json:"address"`
}

type Endpoint struct {
	BaseInfo
	Path string `json:"path"`
}

type BaseInfo struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func HasRole(device Device, role string) bool {
	role = strings.ToLower(role)
	for i := range device.Roles {
		if strings.EqualFold(strings.ToLower(device.Roles[i].ID), role) {
			return true
		}
	}
	return false
}
