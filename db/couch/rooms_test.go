package couch

import (
	"testing"

	"github.com/byuoitav/common/structs"
)

var testRoom = "new_room_a.json"

func TestRoom(t *testing.T) {
	/* should fail
	wipeDatabases()
	t.Run("CreateRoomWithoutBuilding", testCreateRoomWithoutBuilding)

	// setup an initial building to test with
	building := getTestBuilding(t)
	building.ID = "AAA"

	_, err := couch.CreateBuilding(building)
	if err != nil {
		t.Fatalf("failed to create building: %s", err)
	}

	t.Run("CreateRoom", testCreateBuilding)
	wipeDatabase("rooms")
	*/

	t.Run("TestGettingRoom", testGettingRoom)
	//t.Run("GetRoom", testGetBuilding)
	//wipeDatabase("rooms")

	//t.Run("UpdateRoom", testBuildingUpdate)
	//wipeDatabase("rooms")

	//t.Run("DeleteBuilding", testDeleteBuilding)

	//	wipeDatabases()
}

func testGettingRoom(t *testing.T) {
	room, err := couch.getRoom("ITB-1108")
	if err != nil {
		t.Fatalf("error geting test room: %s", err)
	}

	t.Logf("room: %v", room)
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
		t.Fatalf("failed to create building: %s", err)
	}
}

func testGetRoom(t *testing.T) {
}

func getTestRoom(t *testing.T) structs.Room {
	var room structs.Room

	err := unmarshalFromFile(testRoom, &room)
	if err != nil {
		t.Fatalf("failed to unmarshal %s: %s", testRoom, err)
	}

	return room
}
