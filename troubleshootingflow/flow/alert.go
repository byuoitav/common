package flow

import "time"

type Alert struct {
	ID        string
	Timestamp time.Time
	DeviceID  string
	Message   string
}
