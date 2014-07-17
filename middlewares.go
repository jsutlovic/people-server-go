package main

import (
	"fmt"
	"github.com/gocraft/web"
	"log"
	"net/http"
)

func (c *AuthContext) AuthRequired(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	// Get API key
	// Get Username
	// check against database

	_, creds, err := GetAuthHeader(req.Request)
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
