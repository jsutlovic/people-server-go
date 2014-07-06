package main

import (
	"encoding/json"
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
	if !(email_ok || password_ok) {
		status := http.StatusBadRequest
		rw.WriteHeader(status)
		fmt.Fprint(rw, "Email and password are required")
		return
	}

	email := emails[0]
	password := passwords[0]

	log.Print("Got email: ")
	log.Println(email)
	log.Print("Got password: ")
	log.Print(password)

	user := GetUser(email)
	authed := user != nil && user.CheckPassword(password)
	if authed {
		log.Println("Logged in!")
		user_data, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			fmt.Fprint(rw, "Could not convert to JSON")
			return
		}
		fmt.Fprint(rw, string(user_data))
	} else {
		log.Println("Nope.")
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprint(rw, InvalidCredentials)
	}
}
