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
	vals := roomValidationRegex.FindStringSubmatch(r.ID)
	if len(vals) == 0 {
		return errors.New("invalid room: _id must match `([A-z,0-9]{2,})-[A-z,0-9]+`")
	}

	if len(r.Name) == 0 {
		return errors.New("invalid room: missing name.")
	}

	if len(r.Designation) == 0 {
		return errors.New("invalid room: missing designation.")
	}

	if err := r.Configuration.Validate(false); err != nil {
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

func (rc *RoomConfiguration) Validate(deepCheck bool) error {
	if len(rc.ID) == 0 {
		return errors.New("invalid room configuration: missing _id.")
	}

	if deepCheck {
		if len(rc.Evaluators) == 0 {
			return errors.New("invalid room configuration: at least one evaluator is required.")
		}

		for _, evaluator := range rc.Evaluators {
			if err := evaluator.Validate(); err != nil {
				return errors.New(fmt.Sprintf("invalid room configuration: %s", err))
			}
		}
	}

	return nil
}

type Evaluator struct {
	ID          string   `json:"_id"`
	CodeKey     string   `json:"codekey,omitempty"`
	Description string   `json:"description,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	Tags        []string `json:"tags"`
}

func (e *Evaluator) Validate() error {
	if len(e.ID) == 0 {
		return errors.New("invalid evaluator: missing evaluator _id.")
	}

	if len(e.CodeKey) == 0 {
		return errors.New("invalid evaluator: missing codekey")
	}

	// default priority to 1000
	if e.Priority == 0 {
		e.Priority = 1000
	}

	return nil
}
