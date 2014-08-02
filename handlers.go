package main

import (
	"encoding/json"
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
	JsonMalformedError   = "Malformed JSON"
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

type UserCreate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
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

	dec := json.NewDecoder(req.Body)
	var newUser UserCreate
	err := dec.Decode(&newUser)
	if err != nil {
		http.Error(rw, JsonMalformedError, http.StatusBadRequest)
		return
	}

	if newUser.Email == "" || newUser.Password == "" || newUser.Name == "" {
		http.Error(rw, InvalidUserDataError, http.StatusBadRequest)
	}
}
