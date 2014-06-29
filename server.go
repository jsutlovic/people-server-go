package main

import (
	"fmt"
	"github.com/gocraft/web"
	"log"
	"net/http"
)

type Context struct {
}

func (c *Context) AuthMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	// Get API key
	// Get Username
	// check against database
	next(rw, req)
}

func (c *Context) TestMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	next(rw, req)
}

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

func main() {
	rootRouter := web.New(Context{})

	plainRouter := rootRouter.Subrouter(Context{}, "").
		Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware).
		Middleware((*Context).AuthMiddleware).
		Middleware((*Context).TestMiddleware)
	plainRouter.Get("/", (*Context).TestHello)

	authRouter := rootRouter.Subrouter(Context{}, "")
	authRouter.Post("/auth", (*Context).ApiAuth)

	DbInit()

	http.ListenAndServe("0.0.0.0:3000", rootRouter)
}
