package main

import (
	"fmt"
	"github.com/gocraft/web"
	"net/http"
)

// Basic Context available to all handlers
type Context struct {
	DB DbService
}

// Context supplying an authorized user. Used with AuthRequired middleware
type AuthContext struct {
	*Context
	User *User
}

func main() {
	dbService := NewPgDbService("postgres", "user=vagrant dbname=people host=/var/run/postgresql sslmode=disable application_name=people-go")

	rootRouter := web.New(Context{})
	rootRouter.Middleware(web.LoggerMiddleware)
	rootRouter.Middleware(web.ShowErrorsMiddleware)
	rootRouter.Middleware(DbMiddleware(dbService))

	authRouter := rootRouter.Subrouter(Context{}, "")
	authRouter.Post("/auth", (*Context).ApiAuth)

	// CreateUser endpoint cannot require auth
	createUserRouter := rootRouter.Subrouter(Context{}, "")
	createUserRouter.Post("/api/user", (*Context).CreateUserApi)

	// API subrouter for all other API endpoints
	apiRouter := rootRouter.Subrouter(AuthContext{}, "/api")
	apiRouter.Middleware((*AuthContext).AuthRequired)

	// User-related
	apiRouter.Get("/user", (*AuthContext).GetUserApi)

	// Person-related
	apiRouter.Get("/person/:id:\\d+/", (*AuthContext).GetPersonApi)
	apiRouter.Get("/person/", (*AuthContext).GetPersonListApi)

	fmt.Println("Starting server")
	http.ListenAndServe("0.0.0.0:3000", rootRouter)
}
