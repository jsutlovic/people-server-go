package main

import (
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
	dbService := NewPgDbService()

	rootRouter := web.New(Context{})
	rootRouter.Middleware(web.LoggerMiddleware)
	rootRouter.Middleware(web.ShowErrorsMiddleware)
	rootRouter.Middleware(DbMiddleware(dbService))

	plainRouter := rootRouter.Subrouter(AuthContext{}, "")
	plainRouter.Get("/", (*AuthContext).Index)

	authRouter := rootRouter.Subrouter(Context{}, "")
	authRouter.Post("/auth", (*Context).ApiAuth)

	apiRouter := rootRouter.Subrouter(AuthContext{}, "/api")
	apiRouter.Middleware((*AuthContext).AuthRequired)
	apiRouter.Get("/user", (*AuthContext).GetUserApi)

	http.ListenAndServe("0.0.0.0:3000", rootRouter)
}
