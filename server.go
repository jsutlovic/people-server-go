package main

import (
	"errors"
	"fmt"
	"github.com/gocraft/web"
	"net/http"
	"path"
)

type Server struct {
	conf       Config
	rootRouter *web.Router
	routes     []PathRoute
}

type PrefixRouter struct {
	router     *web.Router
	pathPrefix string
}

type methodAction func(r *web.Router, path string, fn interface{}) *web.Router

type httpMethod struct {
	name   string
	action methodAction
}

type PathRoute struct {
	Method  httpMethod
	Path    string
	Handler interface{}
}

func NewServer(conf Config) *Server {
	serv := new(Server)
	serv.conf = conf
	return serv
}

func NewPrefixSubrouter(root *web.Router, prefix string, context interface{}) *PrefixRouter {
	subrouter := root.Subrouter(context, prefix)
	return &PrefixRouter{subrouter, prefix}
}

var (
	httpMethodGet    = httpMethod{"GET", (*web.Router).Get}
	httpMethodPost   = httpMethod{"POST", (*web.Router).Post}
	httpMethodPut    = httpMethod{"PUT", (*web.Router).Put}
	httpMethodPatch  = httpMethod{"PATCH", (*web.Router).Patch}
	httpMethodDelete = httpMethod{"DELETE", (*web.Router).Delete}
)

// Helper method to register a route with the router and server.routes
func (s *Server) registerRoute(router *PrefixRouter, method httpMethod, routePath string, handler interface{}) {
	action := method.action
	action(router.router, routePath, handler)
	s.routes = append(s.routes, PathRoute{method, path.Join(router.pathPrefix, routePath), handler})
}

func (s *Server) setupRoutes() *web.Router {
	rootRouter := web.New(Context{})
	rootRouter.Middleware(web.LoggerMiddleware)
	rootRouter.Middleware(web.ShowErrorsMiddleware)

	// Routers
	authRouter := NewPrefixSubrouter(rootRouter, "", Context{})

	// CreateUser endpoint cannot require auth
	createUserRouter := NewPrefixSubrouter(rootRouter, "", Context{})

	// API subrouter for all other API endpoints
	apiRouter := NewPrefixSubrouter(rootRouter, "/api", AuthContext{})
	apiRouter.router.Middleware((*AuthContext).AuthRequired)

	// Routes
	s.registerRoute(authRouter, httpMethodPost, "/auth", (*Context).ApiAuth)
	s.registerRoute(createUserRouter, httpMethodPost, "/api/user", (*Context).CreateUserApi)

	// User-related
	s.registerRoute(apiRouter, httpMethodGet, "/user", (*AuthContext).GetUserApi)

	// Person-related
	s.registerRoute(apiRouter, httpMethodGet, "/person/:id:\\d+", (*AuthContext).GetPersonApi)
	s.registerRoute(apiRouter, httpMethodGet, "/person", (*AuthContext).GetPersonListApi)
	s.registerRoute(apiRouter, httpMethodPost, "/person", (*AuthContext).CreatePersonApi)

	return rootRouter
}

func (s *Server) Serve() {
	if s.conf == nil {
		panic(errors.New("Config cannot be nil"))
	}
	fmt.Println("Starting server")

	dbService := NewPgDbService(s.conf.DbType(), s.conf.DbCreds())

	s.rootRouter = s.setupRoutes()
	s.rootRouter.Middleware(DbMiddleware(dbService))

	err := http.ListenAndServe(s.conf.ListenAddr(), s.rootRouter)
	if err != nil {
		panic(err)
	}
}
