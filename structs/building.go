package structs

type Building struct {
	ID string `json:"_id"`
	//	Rev         string   `json:"_rev,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
