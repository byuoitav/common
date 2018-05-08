package db

import "testing"

func TestStuff(t *testing.T) {
	db := GetDB()

	buildings, err := db.GetAllBuildings()
	if err != nil {
		t.Logf("error: %s", err)
	}

	t.Logf("buildings: %s", buildings)
}
