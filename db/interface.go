package db

import (
	"log"
	"os"

	"github.com/byuoitav/common/db/couch"
	"github.com/byuoitav/common/structs"
)

type DB interface {
	/* crud functions */
	// building
	CreateBuilding(building structs.Building) (structs.Building, error)
	GetBuilding(id string) (structs.Building, error)
	UpdateBuilding(id string, building structs.Building) (structs.Building, error)
	DeleteBuilding(id string) error

	// room
	CreateRoom(room structs.Room) (structs.Room, error)
	GetRoom(id string) (structs.Room, error)
	UpdateRoom(id string, room structs.Room) (structs.Room, error)
	DeleteRoom(id string) error

	// device
	CreateDevice(device structs.Device) (structs.Device, error)
	GetDevice(id string) (structs.Device, error)
	UpdateDevice(id string, device structs.Device) (structs.Device, error)
	DeleteDevice(id string) error

	// device type
	CreateDeviceType(dt structs.DeviceType) (structs.DeviceType, error)
	GetDeviceType(id string) (structs.DeviceType, error)
	UpdateDeviceType(id string, dt structs.DeviceType) (structs.DeviceType, error)
	DeleteDeviceType(id string) error

	// room configuration
	CreateRoomConfiguration(rc structs.RoomConfiguration) (structs.RoomConfiguration, error)
	GetRoomConfiguration(id string) (structs.RoomConfiguration, error)
	UpdateRoomConfiguration(id string, rc structs.RoomConfiguration) (structs.RoomConfiguration, error)
	DeleteRoomConfiguration(id string) error

	/* bulk functions */
	GetAllBuildings() ([]structs.Building, error)
	GetAllRooms() ([]structs.Room, error)
	GetAllDevices() ([]structs.Device, error)
	GetAllDeviceTypes() ([]structs.DeviceType, error)
	GetAllRoomConfigurations() ([]structs.RoomConfiguration, error)

	/* Specialty functions */
	GetDevicesByRoomAndRole(roomID, roleID string) ([]structs.Device, error)
	GetDevicesByRoleAndType(roleID, typeID string) ([]structs.Device, error)
}

var address string
var username string
var password string

func init() {
	address = os.Getenv("DB_ADDRESS")
	username = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")

	if len(address) == 0 {
		log.Fatalf("DB_ADDRESS is not set. Failing...")
	}
}

func GetDB() DB {
	// TODO add logic to "pick" which db to create
	return couch.NewDB(address, username, password)
}
