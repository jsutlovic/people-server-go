package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"errors"
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

	ct, ctok := rec.HeaderMap["Content-Type"]
	if !assert.True(t, ctok, "No Content-Type header") {
		return
	}
	assert.Equal(t, ct[0], "application/json")
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

	ct, ctok := rec.HeaderMap["Content-Type"]
	if !assert.True(t, ctok, "No Content-Type header") {
		return
	}
	assert.Equal(t, ct[0], "application/json")
}

func TestApiAuthInvalidForm(t *testing.T) {
	var testInvalidParams []url.Values = []url.Values{
		url.Values{"email": []string{"test@example.com"}},
		url.Values{"email": []string{"test@example.com"}, "pw": []string{"adsf"}},
		url.Values{"username": []string{"test@example.com"}, "password": []string{"asdf"}},
		url.Values{"password": []string{"asdf"}},
		url.Values{"pw": []string{"asdf"}},
	}

	for _, params := range testInvalidParams {
		rw, req, rec := mockHandlerParams("POST", "application/x-www-form-urlencoded", params.Encode())

		ac, dbs := mockAuthContext(nil)

		(*AuthContext).ApiAuth(ac, rw, req)

		dbs.Mock.AssertNotCalled(t, "GetUser")
		assert.Equal(t, rec.Code, http.StatusBadRequest)
		assert.Equal(t, rec.Body.String(), ParamsRequired)
	}
}

func TestApiAuthNoUser(t *testing.T) {
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

	dbs.Mock.On("GetUser", user.Email).Return(nil, errors.New("No user"))

	(*AuthContext).ApiAuth(ac, rw, req)

	dbs.Mock.AssertCalled(t, "GetUser", user.Email)
	assert.Equal(t, rec.Code, http.StatusForbidden)
	assert.Equal(t, rec.Body.String(), InvalidCredentials+"\n")
}

func TestApiAuthWrongPassword(t *testing.T) {
	password := "asdf"
	pwhash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		t.Error("Could not create password")
	}

	user := newTestUser()
	user.Pwhash = string(pwhash)

	data := url.Values{}
	data.Add("email", user.Email)
	data.Add("password", "asdf2")

	rw, req, rec := mockHandlerParams("POST", "application/x-www-form-urlencoded", data.Encode())

	ac, dbs := mockAuthContext(user)

	dbs.Mock.On("GetUser", user.Email).Return(user, nil)

	(*AuthContext).ApiAuth(ac, rw, req)

	dbs.Mock.AssertCalled(t, "GetUser", user.Email)
	assert.Equal(t, rec.Code, http.StatusForbidden)
	assert.Equal(t, rec.Body.String(), InvalidCredentials+"\n")
}

func TestApiAuthInactiveUser(t *testing.T) {
	password := "asdf"
	pwhash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		t.Error("Could not create password")
	}

	user := newTestUser()
	user.Pwhash = string(pwhash)
	user.IsActive = false

	data := url.Values{}
	data.Add("email", user.Email)
	data.Add("password", password)

	rw, req, rec := mockHandlerParams("POST", "application/x-www-form-urlencoded", data.Encode())

	ac, dbs := mockAuthContext(user)

	dbs.Mock.On("GetUser", user.Email).Return(user, nil)

	(*AuthContext).ApiAuth(ac, rw, req)

	dbs.Mock.AssertCalled(t, "GetUser", user.Email)
	assert.Equal(t, rec.Code, http.StatusForbidden)
	assert.Equal(t, rec.Body.String(), InactiveUser+"\n")
}

// CreateUserApi allows only JSON data
func TestCreateUserApiJsonOnly(t *testing.T) {
	rw, req, rec := mockHandlerParams("POST", "application/x-www-form-urlencode", "")

	c, _ := mockDbContext(nil)

	(*Context).CreateUserApi(c, rw, req)

	assert.Equal(t, rec.Code, http.StatusBadRequest)
	assert.Equal(t, rec.Body.String(), JsonContentTypeError+"\n")
}

// CreateUserApi does not allow invalid/malformed JSON
func TestCreateUserApiInvalidUserDataJson(t *testing.T) {
	badformats := []map[string]string{
		map[string]string{
			"user": "test@example.com",
		},
		map[string]string{
			"email": "test@example.com",
		},
		map[string]string{
			"email": "test@example.com",
			"name":  "Test User",
		},
		map[string]string{
			"email":    "test@example.com",
			"password": "asdf",
		},
		map[string]string{
			"password": "asdf",
			"name":     "Test User",
		},
		map[string]string{},
		map[string]string{
			"username": "testuser",
		},
	}

	for _, data := range badformats {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			t.Error(err)
		}
		rw, req, rec := mockHandlerParams("POST", "application/json", string(jsonBytes))

		c, _ := mockDbContext(nil)

		(*Context).CreateUserApi(c, rw, req)

		assert.Equal(t, rec.Code, http.StatusBadRequest)
		assert.Equal(t, rec.Body.String(), InvalidUserDataError+"\n")
	}
}

func TestCreateUserApiMalformedJson(t *testing.T) {
	malformed := []string{
		"",
		"asdf",
		"email=test@example.com",
		"email=test@example.com&name=Test+User",
		"email=test@example.com&password=asdf",
		"email:test@example.com",
		`"email":"test@example.com"`,
		`"email":"test@example.com","password":"asdf"`,
		`"email":"test@example.com","name":"Test User"`,
	}

	for _, data := range malformed {
		rw, req, rec := mockHandlerParams("POST", "application/json", data)

		c, _ := mockDbContext(nil)

		(*Context).CreateUserApi(c, rw, req)

		assert.Equal(t, rec.Code, http.StatusBadRequest)
		assert.Equal(t, rec.Body.String(), JsonMalformedError+"\n")
	}
}
