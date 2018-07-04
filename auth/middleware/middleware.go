package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/byuoitav/common/log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"

	cas "gopkg.in/cas.v2"
)

/*
	JWTConfig is the default jwt config object to use in an echo middlware using syntax:
	router.Group("", middlware.JWTWithConfig(auth.JWTConfig))

	To access the information stored in the JWT  use something like:
	claims, ok := context.Get("client").(*jwt.Token).Claims.(jwt.MapClaims)

	Where `claims` will reutrn a key value string of claims validated by the JWT.
*/
var JWTConfig middleware.JWTConfig

func init() {
	if len(os.Getenv("JWT_SIGNING_TOKEN") < 1) {
		log.L.Fatalf("No JWT signing token specified, but using user based authentication. set JWT_SIGNING_TOKEN variable.")
	}
	JWTConfig = middleware.JWTConfig{
		SigningKey:  os.Getenv("JWT_SIGNING_TOKEN"),
		ContextKey:  "client",
		TokenLookup: "header:Authorization",
		AuthScheme:  "Cookie",
	}
}

func AuthenticateUser(next http.Handler) http.Handler {
	url, _ := url.Parse("https://cas.byu.edu/cas")
	c := cas.NewClient(&cas.Options{
		URL: url,
	})

	return c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {

		//check to see if the cookie contains the jwt token
		auth, err := r.Cookie("JWT-TOKEN")
		if err == nil {
			token, err := jwt.Parse(auth, func(T *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != "HS256" {
					return fmt.Errorf("Invalid signing method %v", t.Method.Alg())
				}
				return os.Getenv("JWT_SIGNING_TOKEN")
			})
			if err == nil && token.Valid {
				log.L.Debugf("Adding client claims to context.")
				ctx := context.WithValue(r.Context(), "client", token)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}

		// if they aren't currently authenticated, redirect to the authentication page
		if !cas.IsAuthenticated(r) {
			cas.RedirectToLogin(w, r)
			return
		}

		//otherwise we need to issue a jwt token to the user, as this is the second time they've been here, and they've alreayd authenticated with CAS
	})
}
