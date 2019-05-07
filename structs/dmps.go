package structs

// DMPSList - the list of DMPSes to connect to and pull events
type DMPSList struct {
	ID   string `json:"_id"`
	List []DMPS `json:"list"`
}

// DMPS - a single DMPS to connect to and pull events
type DMPS struct {
	Hostname       string `json:"hostname"`
	Address        string `json:"address"`
	CommandToQuery string `json:"commandToQuery,omitempty"`
	Port           string `json:"port,omitempty"`
}
