package main

import (
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
		{httpMethodPost, "/auth", (*Context).ApiAuth},
		{httpMethodPost, "/api/user", (*Context).CreateUserApi},
		{httpMethodGet, "/api/user", (*AuthContext).GetUserApi},
		{httpMethodGet, "/api/person/:id:\\d+", (*AuthContext).GetPersonApi},
		{httpMethodGet, "/api/person", (*AuthContext).GetPersonListApi},
		{httpMethodPost, "/api/person", (*AuthContext).CreatePersonApi},
	}

	assert.Equal(t, serv.routes, expectedRoutes)
}
