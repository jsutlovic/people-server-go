package main

import (
	"fmt"
	"github.com/gocraft/web"
	"log"
	"net/http"
)

/*
Middleware to require authorization via API key

Checks for the Authorization HTTP header, with Apikey scheme
Extracts credentials and authenticates
If successful, sets a User to the current AuthContext
*/
func (c *AuthContext) AuthRequired(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	_, creds, err := GetAuthHeader(req.Request.Header)
	if err != nil {
		log.Println(err)
		UnauthorizedHeader(rw)
		return
	}

	email, apikey, err := ParseCredentials(creds)
	if err != nil {
		log.Println("Error: could not authenticate")
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := GetUser(email)
	if err != nil {
		log.Println("Error: could not get user from auth")
		log.Println(err)
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprint(rw, "Invalid user")
		return
	}

	if !user.CheckApiKey(apikey) {
		log.Println("Error: API key does not match")
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprint(rw, "Incorrect API key")
		return
	}

	c.User = user

	next(rw, req)
}
