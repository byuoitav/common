package auth

import (
	"fmt"
	"testing"

	"github.com/byuoitav/common/log"
)

func TestAuth(t *testing.T) {
	log.SetLevel("debug")
	test, err := CheckRolesForUser("service", "Ginger", "read-state", "ITB-1101", "room")
	fmt.Printf("Result: %v\nErr:%v", test, err)

}
