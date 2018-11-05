package auth

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/labstack/echo"
)

//CheckAuthForLocalEndpoints assumes that the resource ID being acessed is the room the device is currently in.
func CheckAuthForLocalEndpoints(context echo.Context, role string) (bool, error) {
	ip, _, err := net.SplitHostPort(context.Request().RemoteAddr)
	if err != nil {
		return false, fmt.Errorf("Couldn't parse remote address: %v", err.Error())
	}
	log.L.Debugf("Remote addr: %v", ip)

	if len(os.Getenv("BYPASS_AUTH")) > 0 || ip == "127.0.0.1" || ip == "::1" {
		return true, nil
	}

	accessKeyFromRequest := context.Request().Header.Get("x-av-access-key")
	userFromRequest := context.Request().Header.Get("x-av-user")

	roomID := strings.Split(os.Getenv("SYSTEM_ID"), "-")
	if len(roomID) < 3 {
		return false, fmt.Errorf("couldn't check auth from system id: Invalid systemID")
	}

	return CheckRolesForUser(userFromRequest, accessKeyFromRequest, role, roomID[0]+"-"+roomID[1], "room")
}
