package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetUserApi(t *testing.T) {
	rw, req, rec := mockHandlerParams("GET")

	user := newTestUser()

	ac, _ := mockAuthContext(user)

	(*AuthContext).GetUserApi(ac, rw, req)

	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), Jsonify(user))
}
