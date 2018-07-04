package activedirectory

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var ldapUsername string
var ldapPassword string
var ldapURL string
var userSearch strings.Builder

func init() {
	ldapUsername = os.Getenv("LDAP_USERNAME")
	ldapPassword = os.Getenv("LDAP_PASSWORD")
	ldapURL = os.Getenv("LDAP_URL")

	if len(ldapUsername) == 0 || len(ldapPassword) == 0 || len(ldapURL) == 0 {
		log.Fatalf("LDAP username, password, or URL not set.")
	}

	// build the user search string
	userSearch.WriteString("OU=People")
	split := strings.Split(ldapURL, ":")[0]
	dc := strings.Split(split, ".")

	for _, d := range dc {
		userSearch.WriteString(fmt.Sprintf(",DC=%s", d))
	}
}

func GetGroupsForUser(user string) ([]string, error) {
	var groups []string
	// connect to ldap server
	l, err := ldap.Dial("tcp", ldapURL)
	if err != nil {
		return groups, errors.New(fmt.Sprintf("unable to get groups: %s", err))
	}
	defer l.Close()

	// connect with tls
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return groups, errors.New(fmt.Sprintf("unable to connect to active directory with tls: %s", err))
	}

	// bind with user/pass
	err = l.Bind(ldapUsername, ldapPassword)
	if err != nil {
		return groups, errors.New(fmt.Sprintf("unable to bind username/password to ldap connection: %s", err))
	}

	// build the search request
	searchRequest := ldap.NewSearchRequest(
		userSearch.String(),
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		fmt.Sprintf("(name=%s)", user),
		[]string{"name", "memberOf"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return groups, errors.New(fmt.Sprintf("failed to search ldap: %s", err))
	}

	for _, entry := range sr.Entries {
		if strings.EqualFold(user, entry.GetAttributeValue("name")) {
			tmp := entry.GetAttributeValues("memberOf")
			groups = translateAttributes(tmp)
			break
		}
	}
	return groups, nil
}

func translateAttributes(attributes []string) []string {
	var paths []string

	for _, attribute := range attributes {
		var path strings.Builder

		vals := strings.Split(attribute, ",")
		vals = reverseStringSlice(vals)

		// get the full path (OU/OU/CN) of each group
		for _, val := range vals {
			v := strings.Split(val, "=")

			switch v[0] {
			case "DC":
				break
			case "OU":
				path.WriteString(v[1])
				path.WriteString("/")
			case "CN":
				paths = append(paths, path.String()+v[1])
			}
		}
	}

	return paths
}

func reverseStringSlice(s []string) []string {
	for i := len(s)/2 - 1; i >= 0; i-- {
		opp := len(s) - 1 - i
		s[i], s[opp] = s[opp], s[i]
	}

	return s
}
