package structs

// ScheduleConfig represents the configuration values necessary for the Calendar service to function properly
type ScheduleConfig struct {
	ID              string `json:"_id"`
	Rev             string `json:"_rev"`
	Resource        string `json:"resource"`
	Name            string `json:"displayname"`
	AutoDiscoverURL string `json:"autodiscover-url"`
	AccessType      string `json:"access-type"`
	Image           string `json:"image-url"`
	BookNow         bool   `json:"allowbooknow"`
	ShowHelp        bool   `json:"showhelp"`
	CalendarType    string `json:"calendar-type"`
	CalendarName    string `json:"calendar-name"`
}
