package auth

import (
	"testing"
)

func TestGetGroups(t *testing.T) {
	_, err := GetGroupsForUser("bljoseph")
	if err != nil {
		t.Fatalf("failed to get groups: %s", err)
	}

	//	log.Printf("groups: %s", groups)
}
