package structs

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Device struct {
	ID          string     `json:"_id"`
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	Description string     `json:"description"`
	DisplayName string     `json:"display-name"`
	Type        DeviceType `json:"type"`
	Roles       []Role     `json:"roles"`
	Ports       []Port     `json:"ports"`
	Tags        []string   `json:"tags"`
}

var deviceValidationRegex = regexp.MustCompile(`([A-z,0-9]{2,}-[A-z,0-9]+)-[A-z]+[0-9]+`)

func (d *Device) Validate() error {
	vals := deviceValidationRegex.FindStringSubmatch(d.ID)
	if len(vals) == 0 {
		return errors.New("invalid device: inproper id. must match `([A-z,0-9]{2,}-[A-z,0-9]+)-[A-z]+[0-9]+`")
	}

	if len(d.Name) < 3 {
		return errors.New("invalid device: name must be at least 3 characters long.")
	}

	// validate device type
	if err := d.Type.Validate(); err != nil {
		return errors.New(fmt.Sprintf("invalid device: %s", err))
	}

	// validate roles
	if len(d.Roles) == 0 {
		return errors.New("invalid device: must include at least 1 role.")
	}
	for _, role := range d.Roles {
		if err := role.Validate(); err != nil {
			return errors.New(fmt.Sprintf("invalid device: %s", err))
		}
	}

	// validate ports
	for _, port := range d.Ports {
		if err := port.Validate(); err != nil {
			return errors.New(fmt.Sprintf("invalid device: %s", err))
		}
	}

	return nil
}

type DeviceType struct {
	ID          string       `json:"_id"`
	Description string       `json:"description,omitempty"`
	Ports       []Port       `json:"ports,omitempty"`
	PowerStates []PowerState `json:"power-states,omitempty"`
	Commands    []Command    `json:"commands,omitempty"`
	Tags        []string     `json:"tags"`
}

func (dt *DeviceType) Validate() error {
	if len(dt.ID) == 0 {
		return errors.New("invalid device type: missing id")
	}
	return nil
}

type PowerState struct {
	ID          string   `json:"_id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func (ps *PowerState) Validate() error {
	if len(ps.ID) < 3 {
		return errors.New("invalid power state: id must be at least 2 characters long")
	}
	return nil
}

type Port struct {
	ID                string   `json:"_id"`
	FriendlyName      string   `json:"friendly-name,omitempty"`
	SourceDevice      string   `json:"source-device,omitempty"`
	DestinationDevice string   `json:"destination-device,omitempty"`
	Description       string   `json:"description,omitempty"`
	Tags              []string `json:"tags"`
}

func (p *Port) Validate() error {
	if len(p.ID) < 3 {
		return errors.New("invalid port: id must be at least 2 characters long")
	}
	return nil
}

type Role struct {
	ID          string   `json:"_id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func (r *Role) Validate() error {
	if len(r.ID) < 3 {
		return errors.New("invalid role: id must at least 2 characters long")
	}
	return nil
}

type Command struct {
	ID           string       `json:"_id"`
	Description  string       `json:"description"`
	Microservice Microservice `json:"microservice"`
	Endpoint     Endpoint     `json:"endpoint"`
	Tags         []string     `json:"tags"`
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
