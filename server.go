package main

import (
	"github.com/gocraft/web"
	"net/http"
)

type Context struct {
}

type ApiContext struct {
	*Context
	User *User
}

func main() {
	rootRouter := web.New(Context{})
	rootRouter.Middleware(web.LoggerMiddleware)
	rootRouter.Middleware(web.ShowErrorsMiddleware)

	plainRouter := rootRouter.Subrouter(ApiContext{}, "")
	plainRouter.Get("/", (*ApiContext).TestHello)

	authRouter := rootRouter.Subrouter(Context{}, "")
	authRouter.Post("/auth", (*Context).ApiAuth)

	apiRouter := rootRouter.Subrouter(ApiContext{}, "/api")
	apiRouter.Middleware((*ApiContext).AuthRequired)
	apiRouter.Get("/user", (*ApiContext).GetUserApi)

	dbInit()

	http.ListenAndServe("0.0.0.0:3000", rootRouter)
}
