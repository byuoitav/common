package auth

import (
	"net/http"
	"net/url"

	cas "gopkg.in/cas.v2"
)

func AuthenticateUser(next http.Handler) http.Handler {
	url, _ := url.Parse("https://cas.byu.edu/cas")
	c := cas.NewClient(&cas.Options{
		URL: url,
	})

	return c.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Access-Control-Allow-Origin", "*")

		// do machine checks

		// if they aren't currently authenticated, redirect to the authentication page
		if !cas.IsAuthenticated(r) {
			cas.RedirectToLogin(w, r)
			return
		}

		// get active directory groups

	})
}
