package couch

import (
	"testing"

	"github.com/byuoitav/common/structs"
)

var testBuilding = "new_building.json"

func TestBuilding(t *testing.T) {
	wipeDatabase("buildings")

	t.Run("CreateBuilding", testCreateBuilding)
	wipeDatabase("buildings")

	t.Run("GetBuilding", testGetBuilding)
	wipeDatabase("buildings")

	//t.Run("UpdateBuilding", testBuildingUpdate)
	//wipeDatabase("buildings")

	t.Run("DeleteBuilding", testDeleteBuilding)

	wipeDatabases()
}

func testCreateBuilding(t *testing.T) {
	building := getTestBuilding(t)

	_, err := couch.CreateBuilding(building)
	if err != nil {
		t.Fatalf("failed to create building: %s", err)
	}
}

func testGetBuilding(t *testing.T) {
	// create a building to get
	testCreateBuilding(t)

	// try to get that building
	building := getTestBuilding(t)

	// get the building
	b, err := couch.GetBuilding(building.ID)
	if err != nil {
		t.Fatalf("failed to get building %s: %s", building.ID, err)
	}

	if !isEqual(building, b) {
		t.Fatalf("got a different building than expected... \ngot: %s\nexpected: %s", b, building)
	}
}

func testDeleteBuilding(t *testing.T) {
	// create the building and then get it's id
	testGetBuilding(t)

	// try to delete that building
	building := getTestBuilding(t)

	err := couch.DeleteBuilding(building.ID)
	if err != nil {
		t.Fatalf("failed to delete building %s: %s", building.ID, err)
	}

	// try getting the deleted building. that should give an error.
	_, err = couch.GetBuilding(building.ID)
	if err == nil {
		//	t.Fatalf("get building failed (on id=%s): %s", building.ID, err)

		t.Fatalf("building %s didn't really get deleted, but the DeleteBuilding() acted like it did.", building.ID)
	}
}

func testBuildingUpdate(t *testing.T) {
	testCreateBuilding(t)

	building := getTestBuilding(t)

	// save the oldID to update
	oldID := building.ID

	// modify the building
	building.Name = "updated building name"
	building.ID = "NEWID"
	building.Description = "updated building description"
	building.Tags = []string{"blue", "red", "purple"}

	// update the building
	_, err := couch.UpdateBuilding(oldID, building)
	if err != nil {
		t.Fatalf("failed to update building %s: %s", oldID, err)
	}

	b, err := couch.GetBuilding(building.ID)
	if err != nil {
		t.Fatalf("failed to get updated building %s: %s", building.ID, err)
	}

	if !isEqual(b, building) {
		t.Fatalf("updated building doesn't match expected building.\ngot: %s\nexpected: %s", b, building)
	}
}

func isEqual(b1 structs.Building, b2 structs.Building) bool {
	if b1.ID != b2.ID ||
		b1.Name != b2.Name ||
		b1.Description != b2.Description ||
		len(b1.Tags) != len(b2.Tags) {
		return false
	}

	for i := range b1.Tags {
		if b1.Tags[i] != b2.Tags[i] {
			return false
		}
	}

	return true
}

func getTestBuilding(t *testing.T) structs.Building {
	var building structs.Building

	err := unmarshalFromFile(testBuilding, &building)
	if err != nil {
		t.Fatalf("failed to unmarshal %s: %s", testBuilding, err)
	}

	return building
}

/*
func testBuildingCreateDuplicate(t *testing.T) {

	building := structs.Building{}
	//add a building
	err := structs.UnmarshalFromFile(testDir+"/setup_buildings_a.json", &building)
	if err != nil {
		t.Logf("Error reading in %v: %v", "setup_buildings_a.json", err.Error())
		t.Fail()
	}

	_, err = CreateBuilding(building)
	if err == nil {
		t.Logf("Creation succeeded when should have failed.")
		t.Fail()
	}
}
*/
