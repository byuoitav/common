package auth

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/byuoitav/common/auth/activedirectory"
	"github.com/byuoitav/common/db"
)

func VerifyRoleForUser(user, role string) (bool, error) {
	// get groups that the user is in from active directory
	groups, nerr := activedirectory.GetGroupsForUser(user)
	if nerr != nil {
		return false, errors.New(fmt.Sprintf("failed to get groups for user: %s", nerr))
	}

	// get roles database
	db := db.GetDB()
	auth, err := db.GetAuth()
	if err != nil {
		return false, errors.New(fmt.Sprintf("failed to get roles database: %s", err))
	}

	var groupsWithRole []*regexp.Regexp

	// build a map of all the groups that have this role
	for _, permission := range auth.Permissions {
		for _, r := range permission.Roles {
			if strings.EqualFold(r, role) {
				groupsWithRole = append(groupsWithRole, regexp.MustCompile(permission.Group))
				continue
			}
		}
	}

	// check if one of the groups the user is part of has the role
	for _, groupRegex := range groupsWithRole {
		for _, group := range groups {
			if groupRegex.MatchString(group) {
				return true, nil
			}
		}
	}

	return false, nil
}
