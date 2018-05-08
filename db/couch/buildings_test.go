package couch

import (
	"testing"

	"github.com/byuoitav/configuration-database-microservice/structs"
	"github.com/stretchr/testify/assert"
)

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
