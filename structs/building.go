package structs

import "errors"

type Building struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func (b *Building) Validate() error {
	if len(b.ID) == 0 {
		return errors.New("building must have an id")
	}
	return nil
}
