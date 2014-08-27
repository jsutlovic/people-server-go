package main

import (
	"github.com/gocraft/web"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServeNoConfig(t *testing.T) {
	serv := NewServer(nil)
	assert.Panics(t, func() {
		serv.Serve()
	})
}

func TestSetupRoutes(t *testing.T) {
	conf := newTestConfig()
	serv := NewServer(conf)

	_ = serv.setupRoutes()

	expectedRoutes := []PathRoute{
		PathRoute{web.HttpMethodPost, "/auth", (*Context).ApiAuth},
		PathRoute{web.HttpMethodPost, "/api/user", (*Context).CreateUserApi},
		PathRoute{web.HttpMethodGet, "/api/user", (*AuthContext).GetUserApi},
		PathRoute{web.HttpMethodGet, "/api/person/:id", (*AuthContext).GetPersonApi},
		PathRoute{web.HttpMethodGet, "/api/person", (*AuthContext).GetPersonListApi},
	}

	assert.Equal(t, serv.routes, expectedRoutes)
}
