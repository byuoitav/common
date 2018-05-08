package db

import "testing"

func TestStuff(t *testing.T) {
	db := GetDB()

	rooms, err := db.GetAllRooms()
	if err != nil {
		t.Logf("error: %s", err)
	}

	t.Logf("rooms: %v", rooms)
	t.Fail()
}
