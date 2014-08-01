package main

import (
	"fmt"
	"github.com/gocraft/web"
	"net/http"
)

const (
	InvalidCredentials   = "Invalid credentials"
	ParamsRequired       = "Email and password are required"
	InactiveUser         = "User disabled"
	JsonContentType      = "application/json"
	JsonContentTypeError = "Content-Type is not JSON"
	InvalidUserDataError = "Invalid User data"
)

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
		fmt.Fprint(rw, ParamsRequired)
		return
	}

	email := emails[0]
	password := passwords[0]

	user, err := c.DB.GetUser(email)
	authed := err == nil && user != nil && user.CheckPassword(password)
	if authed {
		if user.IsActive {
			rw.Header().Set("Content-Type", JsonContentType)
			fmt.Fprint(rw, Jsonify(user))
		} else {
			http.Error(rw, "User disabled", http.StatusForbidden)
		}
	} else {
		http.Error(rw, InvalidCredentials, http.StatusForbidden)
	}
}

/*
Handler for the GET User API

Returns a JSON representation of the currently authenticated User
*/
func (c *AuthContext) GetUserApi(rw web.ResponseWriter, req *web.Request) {
	rw.Header().Set("Content-Type", JsonContentType)
	fmt.Fprint(rw, Jsonify(c.User))
}

/*
Handler for the POST User API

Takes a JSON representation of a user and creates an account for it
If a duplicate email is found, return 409 Conflict
Otherwise, if creation is successful return 201 Created
*/
func (c *Context) CreateUserApi(rw web.ResponseWriter, req *web.Request) {
	ct, ctok := req.Header["Content-Type"]
	if !ctok || len(ct) < 1 || (len(ct) >= 1 && ct[0] != "application/json") {
		http.Error(rw, JsonContentTypeError, http.StatusBadRequest)
		return
	}
}
