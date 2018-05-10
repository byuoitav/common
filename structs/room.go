package structs

import (
	"errors"
	"fmt"
	"regexp"
)

type Room struct {
	ID            string            `json:"_id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Configuration RoomConfiguration `json:"configuration"`
	Designation   string            `json:"designation"`
	Devices       []Device          `json:"devices,omitempty"`
	Tags          []string          `json:"tags"`
}

var roomValidationRegex = regexp.MustCompile(`([A-z,0-9]{2,})-[A-z,0-9]+`)

func (r *Room) Validate() error {
	vals := roomValidationRegex.FindAllStringSubmatch(r.ID, 1)
	if len(vals) == 0 {
		return errors.New("invalid room: _id must match `([A-z,0-9]{2,})-[A-z,0-9]+`")
	}

	if len(r.Name) == 0 {
		return errors.New("invalid room: missing name.")
	}

	if len(r.Designation) == 0 {
		return errors.New("invalid room: missing designation.")
	}

	if err := r.Configuration.Validate(); err != nil {
		return errors.New(fmt.Sprintf("invalid room: %s", err))
	}

	return nil
}

type RoomConfiguration struct {
	ID          string      `json:"_id"`
	Evaluators  []Evaluator `json:"evaluators,omitempty"`
	Description string      `json:"description,omitempty"`
	Tags        []string    `json:"tags"`
}

func (rc *RoomConfiguration) Validate() error {
	if len(rc.ID) == 0 {
		return errors.New("invalid room configuration: missing _id.")
	}

	return nil
}

type Evaluator struct {
	ID          string   `json:"_id"`
	CodeKey     string   `json:"code-key,omitempty"`
	Description string   `json:"description,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	Tags        []string `json:"tags"`
}

func (e *Evaluator) Validate() error {
	if len(e.ID) == 0 {
		return errors.New("invalid evaluator: missing evaluator _id.")
	}

	return nil
}
