package main

import (
	"errors"
	"fmt"
	"github.com/gocraft/web"
	"net/http"
	"path"
)

type Config struct {
	DbType     string
	DbCreds    string
	ListenAddr string
}

type Server struct {
	conf       *Config
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
	Method  web.HttpMethod
	Path    string
	Handler interface{}
}

func NewServer(conf *Config) *Server {
	serv := new(Server)
	serv.conf = conf
	return serv
}

func NewPrefixSubrouter(root *web.Router, prefix string, context interface{}) *PrefixRouter {
	subrouter := root.Subrouter(context, prefix)
	return &PrefixRouter{subrouter, prefix}
}

// Utility to map a web.HttpMethod to it's counterpart on web.Router
func methodToRouter(method web.HttpMethod) func(*web.Router, string, interface{}) *web.Router {
	switch method {
	case web.HttpMethodGet:
		return (*web.Router).Get
	case web.HttpMethodPost:
		return (*web.Router).Post
	case web.HttpMethodPut:
		return (*web.Router).Put
	case web.HttpMethodDelete:
		return (*web.Router).Delete
	case web.HttpMethodPatch:
		return (*web.Router).Patch
	}
	return nil
}

// Helper method to register a route with the router and server.routes
func (s *Server) registerRoute(router *PrefixRouter, method web.HttpMethod, routePath string, handler interface{}) {
	action := methodToRouter(method)
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
	s.registerRoute(authRouter, web.HttpMethodPost, "/auth", (*Context).ApiAuth)
	s.registerRoute(createUserRouter, web.HttpMethodPost, "/api/user", (*Context).CreateUserApi)

	// User-related
	s.registerRoute(apiRouter, web.HttpMethodGet, "/user", (*AuthContext).GetUserApi)

	// Person-related
	s.registerRoute(apiRouter, web.HttpMethodGet, "/person/:id:\\d+", (*AuthContext).GetPersonApi)
	s.registerRoute(apiRouter, web.HttpMethodGet, "/person", (*AuthContext).GetPersonListApi)
	s.registerRoute(apiRouter, web.HttpMethodPost, "/person", (*AuthContext).CreatePersonApi)

	return rootRouter
}

func (s *Server) Serve() {
	if s.conf == nil {
		panic(errors.New("Config cannot be nil"))
	}
	fmt.Println("Starting server")

	dbService := NewPgDbService(s.conf.DbType, s.conf.DbCreds)

	s.rootRouter = s.setupRoutes()
	s.rootRouter.Middleware(DbMiddleware(dbService))

	err := http.ListenAndServe(s.conf.ListenAddr, s.rootRouter)
	if err != nil {
		panic(err)
	}
}
