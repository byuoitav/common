package structs

// Person - our authentication struct.
type Person struct {
	ID   string `json:"net-id"`
	Name string `json:"name"`
}

// HasAllPeople compares two list of Persons
func HasAllPeople(oldPeople []Person, newPeople ...Person) bool {
	for i := range oldPeople {
		hasTag := false

		for j := range newPeople {
			if newPeople[j].ID == oldPeople[i].ID && newPeople[j].Name == oldPeople[i].Name {
				hasTag = true
				continue
			}
		}

		if !hasTag {
			return false
		}
	}

	return true
}
