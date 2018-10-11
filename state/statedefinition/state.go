package statedefinition

import "time"

type State struct {
	ID    string      // id of the document to update in elk
	Key   string      // key to update in a entry in the static index
	Time  time.Time   // time the state took effect, must be later than the one saved to be stored.
	Tags  []string    //tags
	Value interface{} // value of key to set in static index
}
