package activedirectory

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	ldap "gopkg.in/ldap.v2"
)

var ldapUsername, ldapPassword, ldapURL, ldapSearchScope string

func init() {
	ldapUsername = os.Getenv("LDAP_USERNAME")
	ldapPassword = os.Getenv("LDAP_PASSWORD")
	ldapURL = os.Getenv("LDAP_URL")
	ldapSearchScope := os.Getenv("LDAP_SEARCH_SCOPE")

	if len(ldapUsername) == 0 || len(ldapPassword) == 0 || len(ldapURL) == 0 || len(ldapSearchScope) == 0 {
		log.L.Fatalf("LDAP username, password, search scope, or URL not set.")
	}
}

func GetGroupsForUser(user string) ([]string, *nerr.E) {
	var groups []string
	// connect to ldap server
	l, err := ldap.Dial("tcp", ldapURL)
	if err != nil {
		return groups, nerr.Translate(err).Addf("Unable to dial LDAP to get groups.")
	}
	defer l.Close()

	// connect with tls
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return groups, nerr.Translate(err).Addf("unable to connect to active directory with tls.")
	}

	// bind with user/pass
	err = l.Bind(ldapUsername, ldapPassword)
	if err != nil {
		return groups, nerr.Translate(err).Addf("unable to bind username/password to ldap connection")
	}

	// build the search request
	searchRequest := ldap.NewSearchRequest(
		ldapSearchScope,
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
		return groups, nerr.Translate(err).Addf("failed to search ldap: %s")
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
