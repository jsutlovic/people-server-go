package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestGetUserApi(t *testing.T) {
	rw, req, rec := mockHandlerParams("GET", "", "")

	user := newTestUser()

	ac, _ := mockAuthContext(user)

	(*AuthContext).GetUserApi(ac, rw, req)

	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), Jsonify(user))
}

func TestApiAuth(t *testing.T) {
	password := "asdf"
	pwhash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		t.Error("Could not create password")
	}

	user := newTestUser()
	user.Pwhash = string(pwhash)

	data := url.Values{}
	data.Add("email", user.Email)
	data.Add("password", password)

	rw, req, rec := mockHandlerParams("POST", "application/x-www-form-urlencoded", data.Encode())

	ac, dbs := mockAuthContext(user)

	dbs.Mock.On("GetUser", user.Email).Return(user, nil)

	(*AuthContext).ApiAuth(ac, rw, req)

	dbs.Mock.AssertCalled(t, "GetUser", user.Email)
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), Jsonify(user))
}
