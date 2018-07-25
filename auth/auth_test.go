package auth

import (
	"os"
	"testing"
)

func TestVerifySuccess(t *testing.T) {
	success, err := VerifyRoleForUser(os.Getenv("NET_ID"), "write")
	if err != nil {
		t.Fatal(err)
	}

	if !success {
		t.Fatalf("failed to verify user, but should have been succcessful.")
	}
}
