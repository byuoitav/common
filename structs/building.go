package structs

import "errors"

type Building struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
}

func (b *Building) Validate() error {
	if len(b.ID) < 2 {
		return errors.New("invalid building: id must be at least 2 characters long")
	}
	return nil
}
