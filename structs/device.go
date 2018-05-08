package structs

import "strings"

type Device struct {
	ID string `json:"_id"`
	//	Rev         string     `json:"_rev,omitempty"`
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	Description string     `json:"description"`
	DisplayName string     `json:"display-name"`
	Type        DeviceType `json:"type"`
	Class       string     `json:"class"`
	Roles       []Role     `json:"roles"`
	Ports       []Port     `json:"ports"`
	Tags        []string   `json:"tags"`
}

type DeviceType struct {
	ID string `json:"_id"`
	//	Rev         string       `json:"_rev,omitempty"`
	Class       string       `json:"class"`
	Description string       `json:"description"`
	Ports       []Port       `json:"ports"`
	PowerStates []PowerState `json:"power-states"`
	Commands    []Command    `json:"commands"`
	Tags        []string     `json:"tags"`
}

type PowerState struct {
	ID          string   `json:"_id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type Port struct {
	ID                string   `json:"_id"`
	FriendlyName      string   `json:"friendly-name"`
	SourceDevice      string   `json:"source-device"`
	DestinationDevice string   `json:"destination-device"`
	Description       string   `json:"description"`
	Tags              []string `json:"tags"`
}

type Role struct {
	ID          string   `json:"_id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type Command struct {
	ID           string       `json:"_id"`
	Description  string       `json:"description"`
	Microservice Microservice `json:"microservice"`
	Endpoint     Endpoint     `json:"endpoint"`
	Tags         []string     `json:"tags"`
}

type DeviceQueryResponse struct {
	Docs     []Device `json:"docs"`
	Bookmark string   `json:"bookmark"`
	Warning  string   `json:"warning"`
}

type Microservice struct {
	ID          string   `json:"_id"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	Tags        []string `json:"tags"`
}

type Endpoint struct {
	ID          string   `json:"_id"`
	Description string   `json:"description"`
	Path        string   `json:"path"`
	Tags        []string `json:"tags"`
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
