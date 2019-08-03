package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/byuoitav/common/log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-cas/cas"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var JWTConfig middleware.JWTConfig
var JWTTTL int
var signingKey []byte

func init() {
	JWTTTL = 604800 //one week
	if len(os.Getenv("JWT_SIGNING_TOKEN")) < 1 {
		log.L.Infof("No JWT signing token defined, autogenerating...")
		signingKey = make([]byte, 64)
		_, err := rand.Read(signingKey)
		if err != nil {
			log.L.Fatalf("Couldn't autogenerate key: %v", err.Error())
		}
		log.L.Infof("Done.")
	} else {
		log.L.Infof("Using provided JWT signing token.")
		signingKey = []byte(os.Getenv("JWT_SIGNING_TOKEN"))
	}
}

/*
	AuthenticateUser uses CAS/JWT authentication to authenticate a user, flow is:

	1. Check for valid, unexpired JWT.
	2. Check to see if request is authenticated with CAS
	3. If no - redirect to CAS login
	4. If valid CAS authentication, issue a JWT token, storing it in a cookie.

	To access the information stored in the JWT  use something like:
	claims, ok := context.Request().Context().Value("client").(*jwt.Token).Claims.(jwt.MapClaims)

	Where `claims` will reutrn a key value string of claims validated by the JWT.

	Included in the context will also be the set groups the user is a part of under the key "user-groups"
	groups, ok := context.Request().Context().Value("user-groups").(map[string]bool)

	the user will also be available in the context
	groups, ok := context.Request().Context().Value("user").(string)
*/
// func AuthenticateCASUser(next http.Handler) http.Handler {
// 	u, _ := url.Parse("https://cas.byu.edu/cas")
// 	c := cas.NewClient(&cas.Options{
// 		URL: u,
// 	})

// 	return c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !cas.IsAuthenticated(r) {
// 			cas.RedirectToLogin(w, r)
// 			return
// 		}

// 		user := cas.Username(r)
// 		log.L.Debugf("Authenticated via CAS: %s", user)

// 		ctx := generateContextNoToken(r, user)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 		return
// 	})
// }

func AuthenticateCASUser(next http.Handler) http.Handler {
	url, _ := url.Parse("https://cas.byu.edu/cas")
	c := cas.NewClient(&cas.Options{
		URL: url,
	})

	return c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		log.L.Debugf("Starting Authentication check")

		//check if bypassauth
		if len(os.Getenv("BYPASS_AUTH")) > 0 {
			log.L.Debugf("Bypassing auth check")
			next.ServeHTTP(w, r)
			return
		}

		//check to see if we've already passed an auth check
		if checkPassedAuthCheck(r) {
			next.ServeHTTP(w, r)
			return
		}

		//check to see if the cookie contains the jwt token
		auth, err := r.Cookie("JWT-TOKEN")
		if err == nil {
			log.L.Debugf("Cookie present, validating...")
			//parse the jwt token out
			token, err := jwt.Parse(auth.Value, func(T *jwt.Token) (interface{}, error) {
				if T.Method.Alg() != "HS256" {
					return "", fmt.Errorf("Invalid signing method %v", T.Method.Alg())
				}
				return []byte(signingKey), nil
			})

			//token was valid, check the expiration time.
			if err == nil && token.Valid {
				log.L.Debugf("Valid token, checking validation")
				exp, ok := token.Claims.(jwt.MapClaims)["exp"]
				if ok {
					//jwt has an expiration time
					t, err := time.Parse(time.RFC3339, exp.(string))

					if err == nil {
						log.L.Debugf("Expiration time %v", t.String())
						//the jwt is still valid
						if !t.Before(time.Now()) {
							//add the claims info to the context and pass the request on
							log.L.Debugf("Valid token. Adding client claims to context.")

							ctx := generateContext(r, token, token.Claims.(jwt.MapClaims)["usr"].(string))
							next.ServeHTTP(w, r.WithContext(ctx))
							return
						}
					}

					log.L.Debugf("Invlaid token, checking CAS")
					//the JWT isn't valid, so we'll fall through to the CAS check.
				}
				log.L.Debug("No expiration time included in token")
			}
			if err != nil {
				log.L.Debugf("Error was not nill when parsing JWT. Error: %v", err.Error())
			} else {
				log.L.Debugf("Token wasn't signed with correct key.")
			}
		}

		log.L.Debugf("Checking CAS")
		// if they aren't currently authenticated, redirect to the authentication page
		if !cas.IsAuthenticated(r) {
			log.L.Info("\n\n\n\n WE'RE GOING TO CAS NOW!!! \n\n\n\n")
			log.L.Debugf("\n\nRedirecting to CAS, not currently authenticated.\n\n")
			c.RedirectToLogin(w, r)
			return
		}
		log.L.Debugf("\n\nAuthenticated via CAS, issuing JWT.\n\n")
		user := cas.Username(r)

		//otherwise we need to issue a jwt token to the user, as this is the second time they've been here, and they've already authenticated with CAS
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Duration(JWTTTL) * time.Second).Format(time.RFC3339),
			"usr": user,
		})

		tokenString, err := token.SignedString([]byte(signingKey))
		if err != nil {
			log.L.Errorf("Couldn't sign JWT: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error creating the user session."))
			return
		}

		log.L.Debugf("JWT generated, passing to next handler")

		//tell the deal to set a cookie
		cook := http.Cookie{
			Name:     "JWT-TOKEN",
			Value:    tokenString,
			HttpOnly: false,
			Secure:   false,
		}
		if os.Getenv("COOKIE_DOMAIN") != "" {
			cook.Domain = os.Getenv("COOKIE_DOMAIN")
		}

		log.L.Debugf("Setting cookie")
		http.SetCookie(w, &cook)

		ctx := generateContext(r, token, user)

		//add values to the context
		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}
func CheckHeaderBasedAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if len(bypassAuth) > 0 {
			log.L.Debugf("Bypassing auth check")
			return next(ctx)
		}

		accessKeyFromRequest := ctx.Request().Header.Get("x-av-access-key")
		userFromRequest := ctx.Request().Header.Get("x-av-user")

		if len(accessKeyFromRequest) == 0 || len(userFromRequest) == 0 {
			//we don't have the access key. Skip to next handler.
			return next(ctx)
		}
		newReqCtx := context.WithValue(ctx.Request().Context(), "user", userFromRequest)
		newReqCtx = context.WithValue(newReqCtx, "access-key", accessKeyFromRequest)
		newReqCtx = context.WithValue(newReqCtx, "passed-auth-check", "true")
		//otherwise we add things to the context
		ctx.SetRequest(ctx.Request().WithContext(newReqCtx))
		return next(ctx)
	}
}

// AuthorizeRequest is an echo middleware function that will check the authorization of a user for a specific resource.
func AuthorizeRequest(role, resourceType string, resourceID func(echo.Context) string) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if len(os.Getenv("BYPASS_AUTH")) > 0 {
				log.L.Debugf("Bypassing auth check")
				return next(ctx)
			}

			if !checkPassedAuthCheck(ctx.Request()) {
				return ctx.String(http.StatusUnauthorized, "unauthorized")
			}

			var user string
			var ok bool
			//get the user (this should at least be filled in).
			userint := ctx.Request().Context().Value("user")
			if userint != nil {
				if user, ok = userint.(string); ok {
					if len(user) == 0 {
						log.L.Errorf("No user on passed auth check...")
						return ctx.String(http.StatusInternalServerError, "Something went wrong")
					}
				}
			}

			var authkey string
			authkeyint := ctx.Request().Context().Value("access-key")
			if userint != nil {
				authkey = accessKey
			} else {
				if authkey, ok = authkeyint.(string); ok {
					if len(authkey) == 0 {
						authkey = accessKey
					}
				} else {
					authkey = accessKey

				}
			}

			resource := resourceID(ctx)

			log.L.Debugf("Resource ID for this auth check is %s (request from %s)", resource, ctx.Path())

			ok, err := CheckRolesForUser(user, accessKey, role, resource, resourceType)
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
