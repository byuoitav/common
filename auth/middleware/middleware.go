package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/byuoitav/common/log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"

	cas "gopkg.in/cas.v2"
)

var JWTConfig middleware.JWTConfig
var JWTTTL int

func init() {
	JWTTTL = 60
	if len(os.Getenv("JWT_SIGNING_TOKEN")) < 1 {
		log.L.Fatalf("No JWT signing token specified, but using user based authentication. set JWT_SIGNING_TOKEN variable.")
	}
	JWTConfig = middleware.JWTConfig{
		SigningKey:  os.Getenv("JWT_SIGNING_TOKEN"),
		ContextKey:  "client",
		TokenLookup: "header:Authorization",
		AuthScheme:  "Cookie",
	}
}

/*
	AuthenticateUser uses CAS/JWT authentication to authenticate a user
	To access the information stored in the JWT  use something like:
	claims, ok := context.Request().Context().Value("client").(*jwt.Token).Claims.(jwt.MapClaims)

	Where `claims` will reutrn a key value string of claims validated by the JWT.

	Included in the context will also be the set groups the user is a part of under the key "user-groups"
	groups, ok := context.Request().Context().Value("user-groups").(map[string]bool)

	the user will also be available in the context
	groups, ok := context.Request().Context().Value("user").(string)
*/
func AuthenticateUser(next http.Handler) http.Handler {
	url, _ := url.Parse("https://cas.byu.edu/cas")
	c := cas.NewClient(&cas.Options{
		URL: url,
	})

	return c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		log.L.Debugf("Starting Authentication check")

		//check to see if the cookie contains the jwt token
		auth, err := r.Cookie("JWT-TOKEN")
		if err == nil {
			log.L.Debugf("Cookie present, validating...")
			//parse the jwt token out
			token, err := jwt.Parse(auth.Value, func(T *jwt.Token) (interface{}, error) {
				if T.Method.Alg() != "HS256" {
					return "", fmt.Errorf("Invalid signing method %v", T.Method.Alg())
				}
				return []byte(os.Getenv("JWT_SIGNING_TOKEN")), nil
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
							ctx := context.WithValue(r.Context(), "client", token)
							next.ServeHTTP(w, r.WithContext(ctx))
							return
						}
					}

					log.L.Debugf("Invlaid token, checking CAS")
					//the JWT isn't valid, so we'll fall through to the CAS check.
				}
			}
		}

		log.L.Debugf("Checking CAS")
		// if they aren't currently authenticated, redirect to the authentication page
		if !cas.IsAuthenticated(r) {
			log.L.Debugf("Redirecting to CAS, not currently authenticated.")
			c.RedirectToLogin(w, r)
			return
		}
		log.L.Debugf("Authenticated via CAS, issuing JWT.")
		user := cas.Username(r)

		//otherwise we need to issue a jwt token to the user, as this is the second time they've been here, and they've already authenticated with CAS
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Add(time.Duration(JWTTTL) * time.Second).Format(time.RFC3339),
			"usr": user,
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_TOKEN")))
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
			Domain:   ".byu.edu",
			HttpOnly: false,
			Secure:   false,
		}

		log.L.Debugf("Setting cookie")
		http.SetCookie(w, &cook)

		next.ServeHTTP(w, r)
		return
	})
}
