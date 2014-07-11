package main

import (
	"github.com/gocraft/web"
	"net/http"
)

type Context struct {
}

type AuthContext struct {
	*Context
	User *User
}

func main() {
	rootRouter := web.New(Context{})
	rootRouter.Middleware(web.LoggerMiddleware)
	rootRouter.Middleware(web.ShowErrorsMiddleware)

	plainRouter := rootRouter.Subrouter(AuthContext{}, "")
	plainRouter.Get("/", (*AuthContext).TestHello)

	authRouter := rootRouter.Subrouter(Context{}, "")
	authRouter.Post("/auth", (*Context).ApiAuth)

	apiRouter := rootRouter.Subrouter(AuthContext{}, "/api")
	apiRouter.Middleware((*AuthContext).AuthRequired)
	apiRouter.Get("/user", (*AuthContext).GetUserApi)

	dbInit()

	http.ListenAndServe("0.0.0.0:3000", rootRouter)
}
