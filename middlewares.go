package main

import (
	"github.com/gocraft/web"
	"net/http"
)

/*
Middleware to hook in the database service

Closes over the DbService provided and simply sets it to each Context on request
*/
func DbMiddleware(s DbService) func(*Context, web.ResponseWriter, *web.Request, web.NextMiddlewareFunc) {
	return func(c *Context, rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
		c.DB = s
		next(rw, req)
	}
}

/*
Middleware to require authorization via API key

Checks for the Authorization HTTP header, with Apikey scheme
Extracts credentials and authenticates
If successful, sets a User to the current AuthContext
*/
func (c *AuthContext) AuthRequired(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	_, creds, err := GetAuthHeader(req.Request.Header)
	if err != nil {
		UnauthorizedHeader(rw)
		return
	}

	email, apikey, err := ParseCredentials(creds)
	if err != nil {
		http.Error(rw, "Invalid authentication params", http.StatusBadRequest)
		return
	}

	user, err := c.DB.GetUser(email)
	if err != nil {
		http.Error(rw, "Invalid user", http.StatusForbidden)
		return
	}

	if !user.CheckApiKey(apikey) {
		http.Error(rw, "Incorrect API key", http.StatusForbidden)
		return
	}

	c.User = user

	next(rw, req)
}
