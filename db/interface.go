package db

import (
	"log"
	"os"
	"sync"

	"github.com/byuoitav/common/db/couch"
	"github.com/byuoitav/common/structs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type DB interface {
	/* crud functions */
	// building
	CreateBuilding(building structs.Building) (structs.Building, error)
	GetBuilding(id string) (structs.Building, error)
	//TODO
	UpdateBuilding(id string, building structs.Building) (structs.Building, error)
	DeleteBuilding(id string) error

	// room
	//TODO
	CreateRoom(room structs.Room) (structs.Room, error)
	//TODO
	GetRoom(id string) (structs.Room, error)
	//TODO
	UpdateRoom(id string, room structs.Room) (structs.Room, error)
	//TODO
	DeleteRoom(id string) error

	// device
	//TODO
	CreateDevice(device structs.Device) (structs.Device, error)
	//TODO
	GetDevice(id string) (structs.Device, error)
	//TODO
	UpdateDevice(id string, device structs.Device) (structs.Device, error)
	//TODO
	DeleteDevice(id string) error

	// device type
	//TODO
	CreateDeviceType(dt structs.DeviceType) (structs.DeviceType, error)
	//TODO
	GetDeviceType(id string) (structs.DeviceType, error)
	//TODO
	UpdateDeviceType(id string, dt structs.DeviceType) (structs.DeviceType, error)
	//TODO
	DeleteDeviceType(id string) error

	// room configuration
	//TODO
	CreateRoomConfiguration(rc structs.RoomConfiguration) (structs.RoomConfiguration, error)
	//TODO
	GetRoomConfiguration(id string) (structs.RoomConfiguration, error)
	//TODO
	UpdateRoomConfiguration(id string, rc structs.RoomConfiguration) (structs.RoomConfiguration, error)
	//TODO
	DeleteRoomConfiguration(id string) error

	/* bulk functions */
	//TODO
	GetAllBuildings() ([]structs.Building, error)
	//TODO
	GetAllRooms() ([]structs.Room, error)
	//TODO
	GetAllDevices() ([]structs.Device, error)
	//TODO
	GetAllDeviceTypes() ([]structs.DeviceType, error)
	//TODO
	GetAllRoomConfigurations() ([]structs.RoomConfiguration, error)
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

func GetDB(logger *zap.SugaredLogger) DB {
	return couch.NewDB(address, username, password, logger)
}

var logger *zap.SugaredLogger
var once sync.Once

func GetDBWithDefaultLogger() DB {
	once.Do(func() {
		CFG := zap.NewDevelopmentConfig()
		CFG.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		CFG.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		l, err := CFG.Build()
		if err != nil {
			log.Panicf("failed to build config for zap logger: %v", err.Error())
		}
		logger = l.Sugar()
		logger.Info("Zap logger started for DB")
	})

	return GetDB(logger)
}
