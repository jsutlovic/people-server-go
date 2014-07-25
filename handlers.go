package main

import (
	"fmt"
	"github.com/gocraft/web"
	"net/http"
)

const InvalidCredentials = "INVALID CREDENTIALS"

func (c *Context) Index(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Welcome to People!")
}

/*
Handler to authenticate a user.

Parses the current request form looking for 'email' and 'password'
fields. Checks the database for the given user and password.

If authentication is successful, return a JSON object representing the
authenticated user. Otherwise returns 403 Forbidden.
*/
func (c *Context) ApiAuth(rw web.ResponseWriter, req *web.Request) {
	req.ParseForm()

	form := req.PostForm

	emails, email_ok := form["email"]
	passwords, password_ok := form["password"]
	if !(email_ok && password_ok) {
		status := http.StatusBadRequest
		rw.WriteHeader(status)
		fmt.Fprint(rw, "Email and password are required")
		return
	}

	email := emails[0]
	password := passwords[0]

	user, err := c.DB.GetUser(email)
	authed := err == nil && user != nil && user.IsActive && user.CheckPassword(password)
	if authed {
		fmt.Fprint(rw, Jsonify(user))
	} else {
		http.Error(rw, InvalidCredentials, http.StatusForbidden)
	}
}

/*
Handler for the GET User API

Returns a JSON representation of the currently authenticated User
*/
func (c *AuthContext) GetUserApi(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, Jsonify(c.User))
}
