package activedirectory

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/mavricknz/ldap"
)

var ldapUsername string
var ldapPassword string
var ldapURL string
var ldapPort int
var userSearch strings.Builder

func init() {
	ldapUsername = os.Getenv("LDAP_USERNAME")
	ldapPassword = os.Getenv("LDAP_PASSWORD")
	ldapURL = os.Getenv("LDAP_URL")
	tempLdapPort := os.Getenv("LDAP_PORT")

	if len(ldapUsername) == 0 || len(ldapPassword) == 0 || len(ldapURL) == 0 || len(tempLdapPort) == 0 {
		log.L.Fatalf("LDAP username, password, port or URL not set.")
	}
	var err error

	ldapPort, err = strconv.Atoi(tempLdapPort)
	if err != nil {
		log.L.Fatalf("Couldn't parse %v for a valid port number.", tempLdapPort)
	}

	// build the user search string
	userSearch.WriteString("OU=People")
	split := strings.Split(ldapURL, ":")[0]
	dc := strings.Split(split, ".")

	for _, d := range dc {
		userSearch.WriteString(fmt.Sprintf(",DC=%s", d))
	}
}

func GetGroupsForUser(user string) ([]string, *nerr.E) {
	groups := []string{}
	conn := ldap.NewLDAPConnection(
		ldapURL,
		uint16(ldapPort),
	)

	err := conn.Connect()
	if err != nil {
		return groups, nerr.Translate(err).Addf("Couldn't connect to ldap server")
	}
	defer conn.Close()
	err = conn.Bind(ldapUsername, ldapPassword)
	if err != nil {
		return groups, nerr.Translate(err).Addf("Couldn't bind to ldap server")

	}
	search := ldap.NewSearchRequest(
		"ou=people,dc=byu,dc=local",
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		fmt.Sprintf("(Name=%s)", user),
		[]string{"Name", "MemberOf"},
		nil,
	)
	res, err := conn.Search(search)
	if err != nil {
		return groups, nerr.Translate(err).Addf("Couldn't search ldap server")
	}

	//log.L.Debugf("%v", res)

	//verify name
	for i := 0; i < len(res.Entries); i++ {
		name := res.Entries[i].GetAttributeValue("Name")
		if name != user {
			continue
		}

		groupsTemp := res.Entries[0].GetAttributeValues("MemberOf")
		groups = translateAttributes(groupsTemp)
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
