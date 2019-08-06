package structs

// LabConfig represents the configuration values neccessary for a Lab Attendance system to function properly
type LabConfig struct {
	ID      string `json:"_id"`
	LabName string `json:"lab_name"`
	LabID   string `json:"lab_id"`
}
