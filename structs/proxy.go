package structs

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/byuoitav/common/nerr"
)

// BuildCommandURL builds the full address for a command based off it's the microservice and endpoint.
// If the device is proxied, the host of the url will be the proxy's address
func (d *Device) BuildCommandURL(commandID string) (string, *nerr.E) {
	findCommand := func(id string, commands []Command) *Command {
		for i := range commands {
			if id == commands[i].ID {
				return &commands[i]
			}
		}
		return nil
	}

	command := findCommand(commandID, d.Type.Commands)
	if command == nil {
		return "", nerr.Createf("error", "unable to build command address: no command with id '%s' found on %s", commandID, d.ID)
	}

	url, err := url.Parse(fmt.Sprintf("%s%s", command.Microservice.Address, command.Endpoint.Path))
	if err != nil {
		return "", nerr.Translate(err).Addf("unable to build command address")
	}
	// match the first command
	for reg, proxy := range d.Proxy {
		r, err := regexp.Compile(reg)
		if err != nil {
			return "", nerr.Translate(err).Addf("unable to build command address")
		}

		if r.MatchString(commandID) {
			// use this proxy
			var host strings.Builder

			oldhost := strings.Split(url.Host, ":")
			newhost := strings.Split(proxy, ":")

			switch len(newhost) {
			case 1: // no port on the proxy url
				host.WriteString(newhost[0])

				// add on the old port if there was one
				if len(oldhost) > 1 {
					host.WriteString(":")
					host.WriteString(oldhost[1])
				}
			case 2: // port present on proxy url
				host.WriteString(newhost[0])
				host.WriteString(":")
				host.WriteString(newhost[1])
			default:
				return "", nerr.Createf("error", "unable to build command address: invalid proxy value '%s' on %s", proxy, d.ID)
			}

			url.Host = host.String()
			break
		}
	}

	return url.String(), nil
}
