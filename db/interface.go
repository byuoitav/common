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
