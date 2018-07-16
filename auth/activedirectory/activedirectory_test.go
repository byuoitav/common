package activedirectory

import (
	"log"
	"os"
	"testing"
)

func TestGetGroups(t *testing.T) {
	groups, err := GetGroupsForUser(os.Getenv("NET_ID"))
	if err != nil {
		t.Fatalf("failed to get groups: %s", err)
	}

	log.Printf("groups: %s", groups)
}
