package main

import (
	"github.com/gocraft/web"
)

func (c *Context) AuthMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	// Get API key
	// Get Username
	// check against database
	next(rw, req)
}

func (c *Context) TestMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	next(rw, req)
}
