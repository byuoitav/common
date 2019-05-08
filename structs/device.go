package structs

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/byuoitav/common/log"
)

// Device - a representation of a device involved in a TEC Pi system.
type Device struct {
	ID          string                 `json:"_id"`
	Name        string                 `json:"name"`
	Address     string                 `json:"address"`
	Description string                 `json:"description"`
	DisplayName string                 `json:"display_name"`
	Type        DeviceType             `json:"type,omitempty"`
	Roles       []Role                 `json:"roles"`
	Ports       []Port                 `json:"ports"`
	Tags        []string               `json:"tags,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`

	// Proxy is a map of regex (matching command id's) to the host:port of the proxy
	Proxy map[string]string `json:"proxy,omitempty"`
}

// DeviceIDValidationRegex is our regular expression for validating the correct naming scheme.
var deviceIDValidationRegex = regexp.MustCompile(`([A-z,0-9]{2,}-[A-z,0-9]+)-[A-z]+[0-9]+`)

// IsDeviceIDValid takes a device id and tells you whether or not it is valid.
func IsDeviceIDValid(id string) bool {
	reg := deviceIDValidationRegex.Copy()

	vals := reg.FindStringSubmatch(id)
	if len(vals) == 0 {
		return false
	}
	return true
}

// Validate checks to see if the device's information is valid or not.
func (d *Device) Validate() error {
	vals := deviceIDValidationRegex.FindStringSubmatch(d.ID)
	if len(vals) == 0 {
		return errors.New("invalid device: inproper id. must match `([A-z,0-9]{2,}-[A-z,0-9]+)-[A-z]+[0-9]+`")
	}

	if len(d.Name) < 2 {
		return errors.New("invalid device: name must be at least 3 characters long")
	}

	// validate device type
	if err := d.Type.Validate(false); err != nil {
		return fmt.Errorf("invalid device: %s", err)
	}

	// validate roles
	if len(d.Roles) == 0 {
		return errors.New("invalid device: must include at least 1 role")
	}
	for _, role := range d.Roles {
		if err := role.Validate(); err != nil {
			return fmt.Errorf("invalid device: %s", err)
		}
	}

	// validate ports
	for _, port := range d.Ports {
		if err := port.Validate(); err != nil {
			return fmt.Errorf("invalid device: %s", err)
		}
	}

	return nil
}

// GetDeviceRoomID returns the room ID portion of the device ID.
func (d *Device) GetDeviceRoomID() string {
	return GetRoomIDFromDevice(d.ID)
}

// GetRoomIDFromDevice .
func GetRoomIDFromDevice(d string) string {
	idParts := strings.Split(d, "-")
	if len(idParts) < 3 {
		log.L.Debugf("invalid ID %v", d)
		return d
	}

	roomID := fmt.Sprintf("%s-%s", idParts[0], idParts[1])
	return roomID
}

// GetCommandByID searches for a specific command and returns it if found.
func (d *Device) GetCommandByID(id string) Command {
	for i := range d.Type.Commands {
		if d.Type.Commands[i].ID == id {
			return d.Type.Commands[i]
		}
	}

	// No command found.
	return Command{}
}

// HasCommand .
func (d *Device) HasCommand(id string) bool {
	for i := range d.Type.Commands {
		if d.Type.Commands[i].ID == id {
			return true
		}
	}

	return false
}

// GetPortFromSrc returns the port going to me from src, and nil if one doesn't exist
func (d *Device) GetPortFromSrc(src string) *Port {
	return d.GetPortFromSrcAndDest(src, d.ID)
}

// GetPortFromSrcAndDest returns the port with a matching src/dest, and nil if one doesn't exist
func (d *Device) GetPortFromSrcAndDest(src, dest string) *Port {
	for i := range d.Ports {
		// log.L.Debugf("checking port %s -> %s", d.Ports[i].SourceDevice, d.Ports[i].DestinationDevice)
		if d.Ports[i].SourceDevice == src && d.Ports[i].DestinationDevice == dest {
			return &d.Ports[i]
		}
	}

	return nil
}

// DeviceType - a representation of a type (or category) of devices.
type DeviceType struct {
	ID          string       `json:"_id"`
	Description string       `json:"description,omitempty"`
	DisplayName string       `json:"display_name,omitempty"`
	Input       bool         `json:"input,omitempty"`
	Output      bool         `json:"output,omitempty"`
	Source      bool         `json:"source,omitempty"`
	Destination bool         `json:"destination,omitempty"`
	Roles       []Role       `json:"roles,omitempty"`
	Ports       []Port       `json:"ports,omitempty"`
	PowerStates []PowerState `json:"power_states,omitempty"`
	Commands    []Command    `json:"commands,omitempty"`
	DefaultName string       `json:"default-name,omitempty"`
	DefaultIcon string       `json:"default-icon,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
}

// Validate checks to make sure that the values of the DeviceType are valid.
func (dt *DeviceType) Validate(deepCheck bool) error {
	if len(dt.ID) == 0 {
		return errors.New("invalid device type: missing id")
	}

	if deepCheck {
		// check all of the ports
		for _, port := range dt.Ports {
			if err := port.Validate(); err != nil {
				return fmt.Errorf("invalid device type: %s", err)
			}
		}

		// check all of the commands
		for _, command := range dt.Commands {
			if err := command.Validate(); err != nil {
				return fmt.Errorf("invalid device type: %s", err)
			}
		}
	}
	return nil
}

// PowerState - a representation of a device's power state.
type PowerState struct {
	ID          string   `json:"_id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
}

// Validate checks to make sure that the PowerState's values are valid.
func (ps *PowerState) Validate() error {
	if len(ps.ID) < 3 {
		return errors.New("invalid power state: id must be at least 3 characters long")
	}
	return nil
}

// Port - a representation of an input/output port on a device.
type Port struct {
	ID                string   `json:"_id"`
	FriendlyName      string   `json:"friendly_name,omitempty"`
	PortType          string   `json:"port_type,omitempty"`
	SourceDevice      string   `json:"source_device,omitempty"`
	DestinationDevice string   `json:"destination_device,omitempty"`
	Description       string   `json:"description,omitempty"`
	Tags              []string `json:"tags,omitempty"`
}

// Validate checks to make sure that the Port's values are valid.
func (p *Port) Validate() error {
	if len(p.ID) < 1 {
		return errors.New("invalid port: id must be at least 3 characters long")
	}
	return nil
}

// Role - a representation of a role that a device plays in the overall system.
type Role struct {
	ID          string   `json:"_id"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Validate checks to make sure that the Role's values are valid.
func (r *Role) Validate() error {
	if len(r.ID) < 3 {
		return errors.New("invalid role: id must at least 3 characters long")
	}
	return nil
}

// Command - a representation of an API command to be executed.
type Command struct {
	ID           string       `json:"_id"`
	Description  string       `json:"description"`
	Microservice Microservice `json:"microservice"`
	Endpoint     Endpoint     `json:"endpoint"`
	Priority     int          `json:"priority"`
	Tags         []string     `json:"tags,omitempty"`
}

// Validate checks to make sure that the Command's values are valid.
func (c *Command) Validate() error {
	if len(c.ID) < 3 {
		return errors.New("invalid command: id must be at least 3 characters long")
	}

	// validate microservice
	if err := c.Microservice.Validate(); err != nil {
		return fmt.Errorf("invalid command: %s", err)
	}

	// validate endpoint
	if err := c.Endpoint.Validate(); err != nil {
		return fmt.Errorf("invalid command: %s", err)
	}
	return nil
}

// BuildCommandAddress builds the full address for a command based off it's the microservice and endpoint
func (c Command) BuildCommandAddress() string {
	if len(c.Microservice.Address) > 0 && len(c.Endpoint.Path) > 0 {
		return fmt.Sprintf("%s%s", c.Microservice.Address, c.Endpoint.Path)
	}

	return ""
}

// Microservice - a representation of a microservice in our API.
type Microservice struct {
	ID          string   `json:"_id"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	Tags        []string `json:"tags,omitempty"`
}

// Validate checks to make sure that the Microservice's values are valid.
func (m *Microservice) Validate() error {
	if len(m.ID) < 3 {
		return errors.New("invalid microservice: id must be at least 3 characters long")
	}

	// validate address
	if _, err := url.ParseRequestURI(m.Address); err != nil {
		return fmt.Errorf("invalid microservice: %s", err)
	}
	return nil
}

// Endpoint - a representation of an API endpoint.
type Endpoint struct {
	ID          string   `json:"_id"`
	Description string   `json:"description"`
	Path        string   `json:"path"`
	Tags        []string `json:"tags,omitempty"`
}

// Validate checks to make sure that the Endpoint's values are valid.
func (e *Endpoint) Validate() error {
	if len(e.ID) < 3 {
		return errors.New("invalid endpoint: id must be at least 3 characters long")
	}

	// validate path
	if _, err := url.ParseRequestURI(e.Path); err != nil {
		return fmt.Errorf("invalid endpoint: %s", err)
	}
	return nil
}

// HasRole checks to see if the given device has the given role.
func HasRole(device Device, role string) bool {
	return device.HasRole(role)
}

// HasRole checks to see if the given device has the given role.
func (d *Device) HasRole(role string) bool {
	role = strings.ToLower(role)
	for i := range d.Roles {
		if strings.EqualFold(strings.ToLower(d.Roles[i].ID), role) {
			return true
		}
	}
	return false
}
