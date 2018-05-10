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
	Roles       []Role     `json:"roles"`
	Ports       []Port     `json:"ports"`
	Tags        []string   `json:"tags"`
}

func (d *Device) Building() string {
	return strings.Split(d.ID, "-")[0]
}

func (d *Device) Room() string {
	return strings.Split(d.ID, "-")[1]
}

type DeviceType struct {
	ID string `json:"_id"`
	//	Rev         string       `json:"_rev,omitempty"`
	Description string       `json:"description,omitempty"`
	Ports       []Port       `json:"ports,omitempty"`
	PowerStates []PowerState `json:"power-states,omitempty"`
	Commands    []Command    `json:"commands,omitempty"`
	Tags        []string     `json:"tags"`
}

type PowerState struct {
	ID          string   `json:"_id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type Port struct {
	ID                string   `json:"_id"`
	FriendlyName      string   `json:"friendly-name,omitempty"`
	SourceDevice      string   `json:"source-device,omitempty"`
	DestinationDevice string   `json:"destination-device,omitempty"`
	Description       string   `json:"description,omitempty"`
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

/*
type DeviceQueryResponse struct {
	Docs     []Device `json:"docs"`
	Bookmark string   `json:"bookmark"`
	Warning  string   `json:"warning"`
}
*/

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
