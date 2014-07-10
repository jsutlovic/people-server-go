package main

import (
	"github.com/gocraft/web"
	"net/http"
)

type Context struct {
}

func main() {
	rootRouter := web.New(Context{})

	plainRouter := rootRouter.Subrouter(Context{}, "").
		Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware).
		Middleware((*Context).AuthMiddleware)
	plainRouter.Get("/", (*Context).TestHello)

	authRouter := rootRouter.Subrouter(Context{}, "")
	authRouter.Post("/api/auth", (*Context).ApiAuth)

	dbInit()

	http.ListenAndServe("0.0.0.0:3000", rootRouter)
}
