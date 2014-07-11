package main

import (
	"fmt"
	"github.com/gocraft/web"
	"log"
	"net/http"
)

const InvalidCredentials = "INVALID CREDENTIALS"

func (c *Context) TestHello(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Hello, world!")
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

	log.Println(email)
	log.Print(password)

	user, err := GetUser(email)
	authed := err == nil && user != nil && user.IsActive && user.CheckPassword(password)
	if authed {
		log.Println("Logged in!")
		fmt.Fprint(rw, Jsonify(user))
	} else {
		log.Printf("Could not log in: %v", err)
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprint(rw, InvalidCredentials)
	}
}

/*
Handler for the GET User API
*/
func (c *ApiContext) GetUserApi(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, Jsonify(c.User))
}
