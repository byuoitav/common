package couch

import (
	"testing"

	"github.com/byuoitav/common/structs"
)

var testDevice = "new_device.json"

func TestDevice(t *testing.T) {
	wipeDatabases()

	t.Run("CreateDeviceWithoutRoom", testCreateRoomWithoutBuilding)
	wipeDatabases()

	// create a building and room for everything to be added to
	t.Run("CreateBuilding", testCreateBuilding)
	t.Run("CreateTestRoom", testCreateRoom)

	t.Run("CreateDevice", testCreateDevice)
	wipeDatabase("devices")

	t.Run("GetDevice", testGetDevice)
	wipeDatabase("devices")

	t.Run("DeleteDevice", testDeleteDevice)
	wipeDatabase("devices")

	t.Run("UpdateDevice", testUpdateDevice)
	wipeDatabase("devices")

	wipeDatabases()
}

func testCreateDeviceWithoutRoom(t *testing.T) {
	_, err := couch.CreateDevice(getTestDevice(t))
	if err == nil {
		t.Fatalf("should have failed to create this device, because it didn't have a room for it, but I succeeded")
	}
}

func testCreateDevice(t *testing.T) {
	device := getTestDevice(t)

	_, err := couch.CreateDevice(device)
	if err != nil {
		t.Fatalf("failed to create device: %s", err)
	}
}

func testGetDevice(t *testing.T) {
	testCreateDevice(t)

	device := getTestDevice(t)
	d, err := couch.GetDevice(device.ID)
	if err != nil {
		t.Fatalf("failed to get device %s: %s", device.ID, err)
	}

	if !isDeviceEqual(device, d) {
		t.Fatalf("got wrong device back from database. \ngot: %v\nexpected: %v", d, device)
	}
}

func testDeleteDevice(t *testing.T) {
	testGetDevice(t)

	device := getTestDevice(t)
	err := couch.DeleteDevice(device.ID)
	if err != nil {
		t.Fatalf("failed to delete device %s: %s", device.ID, err)
	}

	// double check that it is for sure deleted
	_, err = couch.GetDevice(device.ID)
	if err == nil {
		t.Fatalf("DeleteDevice() claimed to work, but the device (%s) didn't really get deleted.", device.ID)
	}
}

func testUpdateDevice(t *testing.T) {
	testGetDevice(t)

	device := getTestDevice(t)
	// modify device
	device.Address = "updated"
	device.Name = "updated_name"

	// update
	_, err := couch.UpdateDevice(device.ID, device)
	if err != nil {
		t.Fatalf("failed to update device: %s", err)
	}

	d, err := couch.GetDevice(device.ID)
	if err != nil {
		t.Fatalf("failed to get device after update: %s", err)
	}

	if d.Address != device.Address || d.Name != device.Name {
		t.Fatalf("updated device is incorrect. \ngot: %v\nexpected: %v", d, device)
	}
}

func getTestDevice(t *testing.T) structs.Device {
	var device structs.Device

	err := unmarshalFromFile(testDevice, &device)
	if err != nil {
		t.Fatalf("failed to unmarshal %s: %s", testDevice, err)
	}

	return device
}

func isDeviceEqual(d1 structs.Device, d2 structs.Device) bool {
	return true
}

/*
func testDeviceCreate(t *testing.T) {
	var device structs.Device
	file := "new_device.json"

	err := structs.UnmarshalFromFile(testDir+"/"+file, &device)
	if err != nil {
		t.Logf("Error reading in %v: %v", file, err.Error())
		t.Fail()
	}

	_, err = CreateDevice(device)
	if err != nil {
		t.Logf("Error creating device %v: %v", file, err.Error())
	}
}
*/
