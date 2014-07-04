package main

import (
	"fmt"
	"github.com/gocraft/web"
	"log"
	"net/http"
)

func (c *Context) TestHello(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Hello, world!")
}

func (c *Context) ApiAuth(rw web.ResponseWriter, req *web.Request) {
	req.ParseForm()

	form := req.PostForm
	//log.Println(form)

	emails, email_ok := form["email"]
	passwords, password_ok := form["password"]
	if !(email_ok || password_ok) {
		status := http.StatusBadRequest
		rw.WriteHeader(status)
		fmt.Fprintf(rw, "Email and password are required")
		return
	}

	email := emails[0]
	password := passwords[0]

	log.Print("Got email: ")
	log.Println(email)
	log.Print("Got password: ")
	log.Print(password)

	authed := CheckPassword(email, password)
	if authed {
		log.Println("Logged in!")
		fmt.Fprintf(rw, "Successfully logged in")
	} else {
		log.Println("Nope.")
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(rw, "Email or password is wrong")
	}
}
