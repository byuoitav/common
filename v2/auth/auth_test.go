package auth

import (
	"fmt"
	"os"
	"testing"
)

func TestAuth(t *testing.T) {
	// log.SetLevel("debug")
	// test, err := CheckRolesForUser("service", "Ginger", "read-state", "ITB-1101", "room")
	// fmt.Printf("Result: %v\nErr:%v", test, err)

	groups, err := GetGroupsForUser(os.Getenv("NET_ID"))
	if err != nil {
		t.Fatalf("failed to get groups: %s", err)
	}

	fmt.Printf("groups: %s\n\n", groups)

	users, err := GetUsersByGroup(os.Getenv("RESPONDER_GROUP"))
	if err != nil {
		t.Fatalf("failed to get the users in a group: %s", err)
	}

	fmt.Printf("\nusers: %s\n\n", users)
}
