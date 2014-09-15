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

func TestMethodToRouter(t *testing.T) {
	tests := []struct {
		in  web.HttpMethod
		out interface{}
	}{
		{
			in:  web.HttpMethodGet,
			out: (*web.Router).Get,
		},
		{
			in:  web.HttpMethodPost,
			out: (*web.Router).Post,
		},
		{
			in:  web.HttpMethodPut,
			out: (*web.Router).Put,
		},
		{
			in:  web.HttpMethodDelete,
			out: (*web.Router).Delete,
		},
		{
			in:  web.HttpMethodPatch,
			out: (*web.Router).Patch,
		},
		{
			in:  web.HttpMethod("INFO"),
			out: (func(*web.Router, string, interface{}) *web.Router)(nil),
		},
	}

	for _, test := range tests {
		assert.Equal(t, methodToRouter(test.in), test.out)
	}
}

func TestSetupRoutes(t *testing.T) {
	conf := newTestConfig()
	serv := NewServer(conf)

	_ = serv.setupRoutes()

	expectedRoutes := []PathRoute{
		{web.HttpMethodPost, "/auth", (*Context).ApiAuth},
		{web.HttpMethodPost, "/api/user", (*Context).CreateUserApi},
		{web.HttpMethodGet, "/api/user", (*AuthContext).GetUserApi},
		{web.HttpMethodGet, "/api/person/:id:\\d+", (*AuthContext).GetPersonApi},
		{web.HttpMethodGet, "/api/person", (*AuthContext).GetPersonListApi},
		{web.HttpMethodPost, "/api/person", (*AuthContext).CreatePersonApi},
	}

	assert.Equal(t, serv.routes, expectedRoutes)
}
