package main

import (
	"errors"
	"fmt"
	"github.com/gocraft/web"
	"net/http"
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

func (s *Server) setupRoutes() *web.Router {
	rootRouter := web.New(Context{})
	rootRouter.Middleware(web.LoggerMiddleware)
	rootRouter.Middleware(web.ShowErrorsMiddleware)

	authRouter := rootRouter.Subrouter(Context{}, "")
	authRouter.Post("/auth", (*Context).ApiAuth)
	s.routes = append(s.routes, PathRoute{web.HttpMethodPost, "/auth", (*Context).ApiAuth})

	// CreateUser endpoint cannot require auth
	createUserRouter := rootRouter.Subrouter(Context{}, "")
	createUserRouter.Post("/api/user", (*Context).CreateUserApi)
	s.routes = append(s.routes, PathRoute{web.HttpMethodPost, "/api/user", (*Context).CreateUserApi})

	// API subrouter for all other API endpoints
	apiRouter := rootRouter.Subrouter(AuthContext{}, "/api")
	apiRouter.Middleware((*AuthContext).AuthRequired)

	// User-related
	apiRouter.Get("/user", (*AuthContext).GetUserApi)
	s.routes = append(s.routes, PathRoute{web.HttpMethodGet, "/api/user", (*AuthContext).GetUserApi})

	// Person-related
	apiRouter.Get("/person/:id:\\d+", (*AuthContext).GetPersonApi)
	s.routes = append(s.routes, PathRoute{web.HttpMethodGet, "/api/person/:id", (*AuthContext).GetPersonApi})
	apiRouter.Get("/person", (*AuthContext).GetPersonListApi)
	s.routes = append(s.routes, PathRoute{web.HttpMethodGet, "/api/person", (*AuthContext).GetPersonListApi})

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
