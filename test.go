package main

import (
	"fmt"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/auth"
)

func main() {
	log.SetLevel("debug")
	test, err := auth.CheckRolesForUser("service", "Ginger", "read-state", "ITB-1101", "room")
	fmt.Printf("Result: %v\nErr:%v", test, err)

}
