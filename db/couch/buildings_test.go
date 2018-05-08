package couch

import (
	"testing"

	"github.com/byuoitav/common/structs"
)

func TestBuilding(t *testing.T) {
	wipeDatabase("buildings")

	t.Run("CreateBuilding", testCreateBuilding)
	wipeDatabase("buildings")

	t.Run("GetBuilding", testGetBuilding)
	wipeDatabase("buildings")

	//	t.Run("UpdateBuilding", testBuildingUpdate)

	t.Run("DeleteBuilding", testDeleteBuilding)

	wipeDatabase("buildings")
}

func testCreateBuilding(t *testing.T) {
	var building structs.Building
	err := unmarshalFromFile("new_building.json", &building)
	if err != nil {
		t.Fatalf("failed to unmarshal %s: %s", "new_building.json", err)
	}

	_, err = db.CreateBuilding(building)
	if err != nil {
		t.Fatalf("failed to create building: %s", err)
	}
}

func testGetBuilding(t *testing.T) {
	// create a building to get
	testCreateBuilding(t)

	// try to get that building
	var building structs.Building
	err := unmarshalFromFile("new_building.json", &building)
	if err != nil {
		t.Fatalf("failed to unmarshal %s: %s", "new_building.json", err)
	}

	// get the building
	b, err := db.GetBuilding(building.ID)
	if err != nil {
		t.Fatalf("failed to get building %s: %s", building.ID, err)
	}

	if !isEqual(building, b) {
		t.Fatalf("got a different building than expected... \nexpected: %s\ngot: %s", building, b)
	}
}

func testDeleteBuilding(t *testing.T) {
	// create the building and then get it's id
	testGetBuilding(t)

	// try to delete that building
	var building structs.Building
	err := unmarshalFromFile("new_building.json", &building)
	if err != nil {
		t.Fatalf("failed to unmarshal %s: %s", "new_building.json", err)
	}

	err = db.DeleteBuilding(building.ID)
	if err != nil {
		t.Fatalf("failed to delete building %s: %s", building.ID, err)
	}

	// try getting the deleted building. that should give an error.
	_, err = db.GetBuilding(building.ID)
	if err == nil {
		t.Fatalf("building %s didn't really get deleted, but the function acted like it did.", building.ID)
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

/*
func TestBuilding(t *testing.T) {
	defer setupDatabase(t)(t)

	t.Run("Building Create", testBuildingCreate)
	t.Run("Building Create Duplicate", testBuildingCreateDuplicate)
	t.Run("Building Update", testBuildingUpdate)
	t.Run("Building Delete", testBuildingDelete)
}

func testBuildingCreate(t *testing.T) {
	building := structs.Building{}
	//add a building
	err := structs.UnmarshalFromFile(testDir+"/new_building.json", &building)
	if err != nil {
		t.Logf("Error reading in %v: %v", "new_building.json", err.Error())
		t.Fail()
	}

	_, err = CreateBuilding(building)
	if err != nil {
		t.Logf("Error creating building %v: %v", "new_building.json", err.Error())
		t.Fail()
	}
}

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

func testBuildingUpdate(t *testing.T) {
	building, err := GetBuildingByID("AAA")
	if err != nil {
		t.Logf("Couldn't get building: %v", err.Error())
		t.Fail()
	}

	currentlen := len(building.Tags)
	building.Tags = append(building.Tags, "pootingmonsterpenguins")
	newDescription := "No! Your great grandaughter had to be a CROSS DRESSER!"
	building.Description = newDescription
	rev := building.Rev

	building.Rev = ""

	//try to fail without rev
	_, err = CreateBuilding(building)
	if err == nil {
		t.Log("Succeeded when it shouldn't have. Failed on rev being null")
		t.FailNow()
	}

	building.Rev = rev

	b, err := CreateBuilding(building)
	if err != nil {
		t.Logf("Failed update: %v", err.Error())
		t.FailNow()
	}
	assert.Equal(t, b.Description, newDescription)
	assert.Equal(t, len(building.Tags), (currentlen + 1))

}

func testBuildingDelete(t *testing.T) {
	err := DeleteBuilding("BBB")
	assert.Nil(t, err)

	err = DeleteBuilding("ZZZ")
	assert.NotNil(t, err)

	//try deleting a building that has a room associated with it

	err = DeleteBuilding("CCC")
	assert.NotNil(t, err)
}
*/
