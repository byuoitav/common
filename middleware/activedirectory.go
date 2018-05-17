package auth

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/ldap.v2"
)

func GetGroupsForUser(user string) ([]string, error) {
	var groups []string

	username := os.Getenv("LDAP_USERNAME")
	password := os.Getenv("LDAP_PASSWORD")

	if len(username) == 0 || len(password) == 0 {
		log.Fatalf("LDAP username or password not set.")
	}

	// connect to ldap server
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "cad3.byu.edu", 3268))
	if err != nil {
		return groups, errors.New(fmt.Sprintf("unable to get groups: %s", err))
	}
	defer l.Close()

	// reconnect with tls
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return groups, errors.New(fmt.Sprintf("unable to connect to active directory with tls: %s", err))
	}

	// bind with user/pass
	err = l.Bind(username, password)
	if err != nil {
		return groups, errors.New(fmt.Sprintf("unable to bind username/password to ldap connection: %s", err))
	}

	// build the search request
	searchRequest := ldap.NewSearchRequest(
		"OU=People,DC=byu,DC=local",
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
	var toReturn []string

	for _, attribute := range attributes {
		log.Printf("attribute: %s", attribute)
		vals := strings.Split(attribute, ",")

		for _, val := range vals {
			v := strings.Split(val, "=")

			switch v[0] {
			case "DC":
				break
			case "OU":
				log.Printf("\t\tadding %s as an OU.", v[1])
			case "CN":
				log.Printf("\t\tadding %s as an CN.", v[1])
			}

			log.Printf("\tsplit: %s", v)
		}

		g := strings.Split(attribute, ",")
		toReturn = append(toReturn, strings.TrimPrefix(g[0], "CN="))
	}

	return toReturn
}

func reverseStringSlice(s []string) {
	for i := len(s)/2 - 1; i >= 0; i-- {
		opp := len(s) - 1 - i
		s[i], s[opp] = s[opp], s[i]
	}
}
