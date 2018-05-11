package couch

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

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

	// main
	dbVersion := &structs.Building{
		ID:          "danny",
		Name:        "test",
		Description: "old",
		Tags:        []string{"hi", "two"},
	}
	log.Printf("dbVersion: %+v", dbVersion)

	db := new(structs.Building)
	*db = *dbVersion

	updatedVersion := structs.Building{
		Name: "new name",
	}
	log.Printf("updatedVersion: %+v", updatedVersion)

	b, err := json.Marshal(updatedVersion)
	if err != nil {
		t.Fatalf("failed to marshal updated version: %s", err)
	}
	log.Printf("bytes: %s", b)

	json.NewDecoder(bytes.NewReader(b)).Decode(&db)
	//json.Unmarshal(b, &dbv)
	log.Printf("after merge: %+v", db)
}

func getRoom(name string) {
	room, err := couch.GetRoom(name)
	if err != nil {
		log.Printf("error: %s", err)
	}

	log.Printf("room: %v", room)
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
