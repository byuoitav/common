package couch

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var couch *CouchDB
var testDir = `./test-data`

func init() {
	couch = &CouchDB{
		// these could be test env vars
		//	address:  os.Getenv("TEST_DB_ADDRESS"),
		//	username: os.Getenv("TEST_DB_USERNAME"),
		//	password: os.Getenv("TEST_DB_PASSWORD"),
		address:  os.Getenv("DB_ADDRESS"),
		username: os.Getenv("DB_USERNAME"),
		password: os.Getenv("DB_PASSWORD"),
	}
}

func unmarshalFromFile(filename string, toFill interface{}) error {
	b, err := ioutil.ReadFile(testDir + "/" + filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, toFill)
	return err
}

func wipeDatabases() {
	couch.MakeRequest("DELETE", "buildings", "", nil, nil)
	couch.MakeRequest("DELETE", "rooms", "", nil, nil)
	couch.MakeRequest("DELETE", "room_configurations", "", nil, nil)
	couch.MakeRequest("DELETE", "devices", "", nil, nil)
	couch.MakeRequest("DELETE", "device_types", "", nil, nil)

	couch.MakeRequest("PUT", "buildings", "", nil, nil)
	couch.MakeRequest("PUT", "rooms", "", nil, nil)
	couch.MakeRequest("PUT", "room_configurations", "", nil, nil)
	couch.MakeRequest("PUT", "devices", "", nil, nil)
	couch.MakeRequest("PUT", "device_types", "", nil, nil)
}

func wipeDatabase(name string) {
	couch.MakeRequest("DELETE", name, "", nil, nil)
	couch.MakeRequest("PUT", name, "", nil, nil)
}

/*

func setupDatabase(t *testing.T) func(t *testing.T) {
	//log.CFG.OutputPaths = []string{}
	//tmp, _ := log.CFG.Build()
	//log.L = tmp.Sugar()

	t.Log("Setting up database for testing")

	//set up our environment variables
	//oldCouchAddress := os.Getenv("COUCH_ADDRESS")
	//oldCouchUsername := os.Getenv("COUCH_USERNAME")
	//oldCouchPassword := os.Getenv("COUCH_PASSWORD")
	//oldLoggingLocation := os.Getenv("LOGGING_FILE_LOCATION")

	//os.Setenv("COUCH_ADDRESS", os.Getenv("COUCH_TESTING_ADDRESS"))
	//os.Setenv("COUCH_USERNAME", os.Getenv("COUCH_TESTING_USERNAME"))
	//os.Setenv("COUCH_PASSWORD", os.Getenv("COUCH_TESTING_PASSWORD"))
	//os.Setenv("LOGGING_FILE_LOCATION", os.Getenv("TEST_LOGGING_FILE_LOCATION"))

	//now we go and set up the database

	//find all of the setup files to be read in

	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		msg := fmt.Sprintf("Couldn't read the database setup director: %v", err.Error())
		t.Log(msg)
		t.FailNow()
	}

	setupScriptRegex := regexp.MustCompile(`setup_([A-Z,a-z]+)`)

	//wipe out the current databases.
	wipeDatabases()

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		//check it for a setup
		matches := setupScriptRegex.FindStringSubmatch(f.Name())
		if len(matches) == 0 {
			continue
		}
		t.Logf("Reading in %v", f.Name())

		switch matches[1] {
		case "buildings":
			building := structs.Building{}
			//add a building
			err := structs.UnmarshalFromFile(testDir+"/"+f.Name(), &building)
			if err != nil {
				t.Logf("couldn't set up database. Error reading in %v: %v", f.Name(), err.Error())
				t.FailNow()
			}

			_, err = CreateBuilding(building)
			if err != nil {
				t.Logf("couldn't set up database. Error creating building %v: %v", f.Name(), err.Error())
				t.FailNow()
			}
		case "devicetypes":
			dt := structs.DeviceType{}
			//add a building
			err := structs.UnmarshalFromFile(testDir+"/"+f.Name(), &dt)
			if err != nil {
				t.Logf("couldn't set up database. Error reading in %v: %v", f.Name(), err.Error())
				t.FailNow()
			}

			_, err = CreateDeviceType(dt)
			if err != nil {
				t.Logf("couldn't set up database. Error creating devicetype %v: %v", f.Name(), err.Error())
				t.FailNow()
			}
		case "rooms":
			dt := structs.Room{}
			//add a building
			err := structs.UnmarshalFromFile(testDir+"/"+f.Name(), &dt)
			if err != nil {
				t.Logf("couldn't set up database. Error reading in %v: %v", f.Name(), err.Error())
				t.FailNow()
			}

			_, err = CreateRoom(dt)
			if err != nil {
				t.Logf("couldn't set up database. Error creating room %v: %v", f.Name(), err.Error())
				t.FailNow()
			}
		case "devices":
			dt := structs.Device{}
			//add a building
			err := structs.UnmarshalFromFile(testDir+"/"+f.Name(), &dt)
			if err != nil {
				t.Logf("couldn't set up database. Error reading in %v: %v", f.Name(), err.Error())
				t.FailNow()
			}

			_, err = CreateDevice(dt)
			if err != nil {
				t.Logf("couldn't set up database. Error creating device %v: %v", f.Name(), err.Error())
				t.FailNow()
			}
		case "roomconfigs":
			dt := structs.RoomConfiguration{}
			//add a building
			err := structs.UnmarshalFromFile(testDir+"/"+f.Name(), &dt)
			if err != nil {
				t.Logf("couldn't set up database. Error reading in %v: %v", f.Name(), err.Error())
				t.FailNow()
			}

			_, err = CreateRoomConfiguration(dt)
			if err != nil {
				t.Logf("couldn't set up database. Error creating roomconfiguration %v: %v", f.Name(), err.Error())
				t.FailNow()
			}
		}
	}

	return func(tarp *testing.T) {
//		os.Setenv("COUCH_ADDRESS", oldCouchAddress)
//		os.Setenv("COUCH_USERNAME", oldCouchUsername)
//		os.Setenv("COUCH_PASSWORD", oldCouchPassword)
//		os.Setenv("LOGGING_FILE_LOCATION", oldLoggingLocation)
	}
}

func wipeDatabases() {
	MakeRequest("DELETE", "buildings", "", nil, nil)
	MakeRequest("DELETE", "rooms", "", nil, nil)
	MakeRequest("DELETE", "room_configurations", "", nil, nil)
	MakeRequest("DELETE", "devices", "", nil, nil)
	MakeRequest("DELETE", "device_types", "", nil, nil)

	MakeRequest("PUT", "buildings", "", nil, nil)
	MakeRequest("PUT", "rooms", "", nil, nil)
	MakeRequest("PUT", "room_configurations", "", nil, nil)
	MakeRequest("PUT", "devices", "", nil, nil)
	MakeRequest("PUT", "device_types", "", nil, nil)
}
*/
