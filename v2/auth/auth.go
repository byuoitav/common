package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/endpoint-authorization-controller/base"
	"github.com/labstack/echo"
)

var internalPermissionCache = make(map[string][]base.Response)

//bypass - :0
var bypassAuth = os.Getenv("BYPASS_AUTH")

//for use on the endpoint
var authURL = os.Getenv("ENDPOINT_AUTHORIZATION_URL")

//for use on the client
var accessKey = os.Getenv("ENDPOINT_ACCESS_KEY")
var endpointUserName = os.Getenv("ENDPOINT_USERNAME")

//AddAuthToRequest to add authorization headers to a request.  For use in a client
func AddAuthToRequest(request *http.Request) {
	request.Header.Add("x-av-access-key", accessKey)
	request.Header.Add("x-av-user", endpointUserName)
}

//AddAuthToRequestForUser to add authorization headers to a request.  For use in a client
func AddAuthToRequestForUser(request *http.Request, userName string) {
	request.Header.Add("x-av-access-key", accessKey)
	request.Header.Add("x-av-user", userName)
}

//CheckRolesForReceivedRequest to check authorization of a user for a specific resource.  For use in an endpoint receiving requests
func CheckRolesForReceivedRequest(context echo.Context, role string, resourceID string, resourceType string) (bool, error) {
	accessKeyFromRequest := context.Request().Header.Get("x-av-access-key")
	userFromRequest := context.Request().Header.Get("x-av-user")
	return CheckRolesForUser(userFromRequest, accessKeyFromRequest, role, resourceID, resourceType)
}

//CheckRolesForUser to check authorization of a user for a specific resource.  For use in an endpoint receiving requests
func CheckRolesForUser(user string, accessKey string, role string, resourceID string, resourceType string) (bool, error) {
	if len(bypassAuth) > 0 {
		return true, nil
	}

	cacheTest := checkCacheForAuth(user, accessKey, role, resourceID, resourceType)

	if cacheTest {
		return true, nil
	}

	//not cached (or no match), so go ahead and make the call
	var authRequestBody base.Request
	authRequestBody.AccessKey = accessKey
	authRequestBody.UserInformation.ID = user
	authRequestBody.UserInformation.AuthMethod = "CAS"
	authRequestBody.UserInformation.ResourceType = resourceType
	authRequestBody.UserInformation.ResourceID = resourceID

	var authResponse base.Response

	if len(authURL) == 0 {
		return false, fmt.Errorf("No ENDPOINT_AUTHORIZATION_URL environment variable set")
	}

	log.L.Debugf("Creating auth request")
	authRequest, err := jsonhttp.CreateRequest("POST", authURL+"/authorize", authRequestBody, nil)

	if err != nil {
		log.L.Debugf("Error creating auth request %v", err.Error())
		return false, err
	}

	err = jsonhttp.ExecuteRequest(authRequest, &authResponse, 3)

	if err != nil {
		log.L.Debugf("Error executing auth request %v", err.Error())
		return false, err
	}

	log.L.Debugf("Auth response received: ", authResponse)

	//add to cache
	internalPermissionCache[user] = append(internalPermissionCache[user], authResponse)

	//recheck cache
	return checkCacheForAuth(user, accessKey, role, resourceID, resourceType), nil
}

func checkCacheForAuth(user string, accessKey string, role string, resourceID string, resourceType string) bool {
	if responseArray, ok := internalPermissionCache[user]; ok {
		//this user has something cached - check there first
		log.L.Debugf("User %v has %v cached permissions", user, len(responseArray))

		for i := 0; i < len(responseArray); i++ {
			var oneResponse = responseArray[i]
			if time.Now().After(oneResponse.TTL) {
				//remove this response
				log.L.Debugf("Response %v has expired as of %v, removing", oneResponse, oneResponse.TTL)
				responseArray = append(responseArray[:i], responseArray[i+1:]...)
				i--
			} else {
				//check and see if this one matches
				log.L.Debugf("Response %v is still valid", oneResponse)
				if thisResourcePermission, ok := responseArray[i].Permissions[resourceID]; ok {
					//check that the roles contains the target role
					for _, oneRole := range thisResourcePermission {
						if oneRole == role {
							return true
						}
					}
				}
			}
		}
	}

	return false
}
