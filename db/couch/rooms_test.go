package couch

import (
	"testing"

	"github.com/byuoitav/common/log"

	"github.com/byuoitav/common/structs"
)

var testRoom = "new_room_a.json"

func TestRoom(t *testing.T) {
	/*
		//	wipeDatabases()
		t.Run("CreateRoomWithoutBuilding", testCreateRoomWithoutBuilding)

		// setup an initial building to test with
		building := getTestBuilding(t)
		building.ID = "AAA"

		_, err := couch.CreateBuilding(building)
		if err != nil {
			t.Fatalf("failed to create building: %s", err)
		}

		//	t.Run("CreateRoom", testCreateBuilding)
		//	wipeDatabase("rooms")

		//t.Run("GetRoom", testGetBuilding)
		//wipeDatabase("rooms")

		//t.Run("UpdateRoom", testBuildingUpdate)
		//wipeDatabase("rooms")

		//t.Run("DeleteBuilding", testDeleteBuilding)

		//	wipeDatabases()
	*/

	/* ROUTER
	router := echo.New()

	router.GET("/room/:name", func(context echo.Context) error {
		roomname := context.Param("name")
		room, err := couch.GetRoom(roomname)
		if err != nil {
			return context.JSON(400, err)
		}

		return context.JSON(200, room)
	})

	router.GET("/building/:name", func(context echo.Context) error {
		name := context.Param("name")
		building, err := couch.GetBuilding(name)
		if err != nil {
			return context.JSON(400, err)
		}

		return context.JSON(200, building)
	})

	router.Start(":9999")
	*/

	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 1000

	// var rooms []structs.Room
	//	couch.ExecuteQuery(query, ROOMS, rooms)
}

func getRoom(name string) {
	room, err := couch.GetRoom(name)
	if err != nil {
		log.L.Infof("error: %s", err)
	}

	log.L.Infof("room: %v", room)
}

func testCreateRoomWithoutBuilding(t *testing.T) {
	room := getTestRoom(t)

	_, err := couch.CreateRoom(room)
	if err == nil {
		t.Fatalf("successfully created room when I shouldn't have (there was no building matching the room)")
	}
}

func testCreateRoom(t *testing.T) {
	room := getTestRoom(t)

	_, err := couch.CreateRoom(room)
	if err != nil {
		t.Fatalf("failed to create room: %s", err)
	}
}

func getTestRoom(t *testing.T) structs.Room {
	var room structs.Room

	err := unmarshalFromFile(testRoom, &room)
	if err != nil {
		t.Fatalf("failed to unmarshal %s: %s", testRoom, err)
	}

	return room
}
