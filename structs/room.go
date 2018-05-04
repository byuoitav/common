package structs

type Room struct {
	ID            string            `json:"_id"`
	Rev           string            `json:"_rev,omitempty"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Tags          []string          `json:"tags"`
	Configuration RoomConfiguration `json:"configuration"`
	Designation   string            `json:"designation"`
	Devices       []Device          `json:"devices"`
}

type RoomConfiguration struct {
	ID          string      `json:"_id"`
	Rev         string      `json:"_rev,omitempty"`
	Name        string      `json:"name,omitempty"`
	Evaluators  []Evaluator `json:"evaluators,omitempty"`
	Description string      `json:"description,omitempty"`
}

type Evaluator struct {
	ID          string `json:"_id"`
	CodeKey     string `json:"code-key,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty"`
}

type RoomQueryResponse struct {
	Docs     []Room `json:"docs"`
	Bookmark string `json:"bookmark"`
	Warning  string `json:"warning"`
}
