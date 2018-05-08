package db

import (
	"github.com/byuoitav/common/db/couch"
	"github.com/byuoitav/configuration-database-microservice/structs"
)

type DB interface {
	/* crud functions */
	// building
	CreateBuilding(building structs.Building) error
	GetBuilding(id string) (structs.Building, error)
	UpdateBuilding(id string, building structs.Building) (structs.Building, error)
	DeleteBuilding(id string) error

	// room
	CreateRoom(room structs.Room) error
	GetRoom(id string) (structs.Room, error)
	UpdateRoom(id string, room structs.Room) (structs.Room, error)
	DeleteRoom(id string) error

	// device
	CreateDevice(device structs.Device) error
	GetDevice(id string) (structs.Device, error)
	UpdateDevice(id string, device structs.Device) (structs.Device, error)
	DeleteDevice(id string) error

	// device type
	CreateDeviceType(dt structs.DeviceType) error
	GetDeviceType(id string) (structs.DeviceType, error)
	UpdateDeviceType(id string, dt structs.DeviceType) (structs.DeviceType, error)
	DeleteDeviceType(id string) error

	// room configuration
	CreateRoomConfiguration(rc structs.RoomConfiguration) error
	GetRoomConfiguration(id string) (structs.RoomConfiguration, error)
	UpdateRoomConfiguration(id string, rc structs.RoomConfiguration) (structs.RoomConfiguration, error)
	DeleteRoomConfiguration(id string) error
}

func GetDB() DB {
	return couch.CouchDB{}
}
