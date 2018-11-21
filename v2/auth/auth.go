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

// AuthorizeRequest is an echo middleware function that will check the authorization of a user for a specific resource.
func AuthorizeRequest(role, resourceType string, resourceID func(echo.Context) string) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if len(bypassAuth) > 0 {
				log.L.Debugf("Bypassing auth check")
				return next(ctx)
			}
			log.L.Debugf("Checking auth for endpoint %s, which requires the role '%s' for resource type '%s'. This request is from %s", ctx.Path(), role, resourceType, ctx.RealIP())

			accessKeyFromRequest := ctx.Request().Header.Get("x-av-access-key")
			userFromRequest := ctx.Request().Header.Get("x-av-user")

			if len(accessKeyFromRequest) == 0 || len(userFromRequest) == 0 {
				return ctx.String(http.StatusBadRequest, "must include 'x-av-access-key' and 'x-av-user' headers")
			}

			resource := resourceID(ctx)
			if len(resource) == 0 {
				return ctx.String(http.StatusInternalServerError, "unable to get resource from request")
			}

			log.L.Debugf("Resource ID for this auth check is %s (request from %s)", resource, ctx.Path())

			ok, err := CheckRolesForUser(userFromRequest, accessKeyFromRequest, role, resource, resourceType)
			if err != nil {
				return ctx.String(http.StatusInternalServerError, fmt.Sprintf("unable to authorize request: %s", err))
			}

			if !ok {
				return ctx.String(http.StatusUnauthorized, "Not authorized")
			}

			return next(ctx)
		}
	}
}

// LookupResourceFromAddress uses the ":address" parameter from the endpoint and returns the resourceID requested.
func LookupResourceFromAddress(ctx echo.Context) string {
	// TODO this should rip out :address and either use the hostname
	// or do a reverse DNS lookup to decide what it's hostname is,
	// and then return the resourceID of the address requested
	// should be used on device communication microservices (sony-control, etc)
	return ""
}

//CheckRolesForUser to check authorization of a user for a specific resource.  For use in an endpoint receiving requests
func CheckRolesForUser(user string, accessKey string, role string, resourceID string, resourceType string) (bool, error) {
	cacheTest := checkCacheForAuth(user, accessKey, role, resourceID, resourceType)
	if cacheTest {
		return true, nil
	}

	//not cached (or no match), so go ahead and make the call
	var authRequestBody base.Request
	authRequestBody.AccessKey = accessKey
	authRequestBody.UserInformation.ID = user
	authRequestBody.UserInformation.AuthMethod = "cas"
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
