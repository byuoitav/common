package structs

// Auth - our authentication struct.
type Auth struct {
	ID          string   `json:"_id"`
	Roles       []string `json:"roles"`
	Permissions []struct {
		Group string   `json:"group"`
		Roles []string `json:"roles"`
	} `json:"permissions"`
}
