package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/endpoint-authorization-controller/base"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var internalPermissionCache = make(map[string]base.Response)

var cacheLock sync.RWMutex

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

// LookupResourceFromAddress uses the ":address" parameter from the endpoint and returns the resourceID requested.
func LookupResourceFromAddress(ctx echo.Context) string {
	// TODO this should rip out :address and either use the hostname
	// or do a reverse DNS lookup to decide what it's hostname is,
	// and then return the resourceID of the address requested
	// should be used on device communication microservices (sony-control, etc)
	return "all"
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
	cacheLock.Lock()
	internalPermissionCache[user] = authResponse
	cacheLock.Unlock()

	//recheck cache
	return checkCacheForAuth(user, accessKey, role, resourceID, resourceType), nil
}

func checkCacheForAuth(user string, accessKey string, role string, resourceID string, resourceType string) bool {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	if oneResponse, ok := internalPermissionCache[user]; ok {
		//this user has something cached - check there first
		log.L.Debugf("User %v has cached permissions", user)

		if time.Now().After(oneResponse.TTL) {
			delete(internalPermissionCache, user)
		} else {
			//check and see if this one matches
			log.L.Debugf("Response %v is still valid", oneResponse)
			if thisResourcePermission, ok := oneResponse.Permissions[resourceID]; ok {
				//check that the roles contains the target role
				for _, oneRole := range thisResourcePermission {
					if oneRole == role {
						return true
					}
				}
			}
		}
	}

	return false
}

func checkPassedAuthCheck(r *http.Request) bool {

	pass := r.Context().Value("passed-auth-check")
	if pass != nil {
		if v, ok := pass.(string); ok {
			if v == "true" {
				log.L.Debugf("Pre-passed auth check, passing CAS check")
				return true
			}
		}
	}

	return false
}

//given the user and client token, we'll get all the groups that the user is a part of and include that in the context.
func generateContext(r *http.Request, clientToken *jwt.Token, username string) context.Context {
	ctx := context.WithValue(r.Context(), "client", clientToken)
	ctx = context.WithValue(ctx, "user", username)
	ctx = context.WithValue(ctx, "passed-auth-check", "true")

	//Get the user groups
	return ctx
}
