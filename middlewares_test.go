package main

import (
	"errors"
	"fmt"
	"github.com/gocraft/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockNext struct {
	mock.Mock
}

func (m *MockNext) Next(rw web.ResponseWriter, req *web.Request) {
	m.Mock.Called(rw, req)
	return
}

type MockDbService struct {
	mock.Mock
	*User
}

func (m *MockDbService) GetUser(email string) (user *User, err error) {
	args := m.Mock.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*User), nil
	}
	return nil, args.Error(1)
}

func mockMiddlewareParams() (web.ResponseWriter, *web.Request, *MockNext, *httptest.ResponseRecorder) {
	// Build the ResponseRecorder
	recorder := httptest.NewRecorder()
	rw := new(web.AppResponseWriter)
	rw.ResponseWriter = recorder

	// Build the request
	fakeRequest, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		panic(err)
	}
	req := web.Request{}
	req.Request = fakeRequest

	// Mock a NextMiddlewareFunc
	next := new(MockNext)

	// Setup expecations for Next
	next.Mock.On("Next", rw, &req).Return()

	return rw, &req, next, recorder
}

func mockDbContext(user *User) (Context, *MockDbService) {
	// Create our mock database service to serve our fake user
	dbs := new(MockDbService)
	dbs.User = user

	// Create Context and set the DB
	c := Context{}
	c.DB = dbs

	return c, dbs
}

func TestDbMiddleware(t *testing.T) {
	// Get our mock service
	dbs := new(MockDbService)

	// Build some basic middleware objects
	rw, req, next, _ := mockMiddlewareParams()

	// Create a Context
	c := new(Context)

	// Build the middleware (closure)
	dbMiddleware := DbMiddleware(dbs)

	// Call the middleware
	dbMiddleware(c, rw, req, next.Next)

	// Assertions
	next.Mock.AssertCalled(t, "Next", rw, req)
	assert.Equal(t, c.DB, dbs)
}

func TestAuthRequiredAuthorizesValid(t *testing.T) {
	user := newTestUser()

	// Build our basic middleware objects
	rw, req, next, rec := mockMiddlewareParams()

	// Add headers to request
	req.Request.Header.Add("Authorization", fmt.Sprintf("Apikey %s:%s", user.Email, user.ApiKey))

	// Setup contexts
	c, dbs := mockDbContext(&user)
	dbs.Mock.On("GetUser", user.Email).Return(&user, nil)

	ac := new(AuthContext)
	ac.Context = c

	// Call the middleware
	(*AuthContext).AuthRequired(ac, rw, req, next.Next)

	// Assertions:

	// Next was called
	next.Mock.AssertCalled(t, "Next", rw, req)
	// GetUser was called
	dbs.Mock.AssertCalled(t, "GetUser", user.Email)
	// User is set to the AuthContext
	assert.Equal(t, ac.User, &user)
	// Nothing was written to the responsewriter
	assert.Equal(t, rec.Body.String(), "")
}

func TestAuthRequiredNoHeader(t *testing.T) {
	user := newTestUser()

	rw, req, next, rec := mockMiddlewareParams()

	c, dbs := mockDbContext(&user)
	dbs.Mock.On("GetUser", user.Email).Return(&user, nil)

	ac := new(AuthContext)
	ac.Context = c

	(*AuthContext).AuthRequired(ac, rw, req, next.Next)

	next.Mock.AssertNotCalled(t, "Next", rw, req)
	dbs.Mock.AssertNotCalled(t, "GetUser", user.Email)
	assert.Nil(t, ac.User)
	assert.Equal(t, rec.Code, http.StatusUnauthorized)
	assert.Equal(t, rec.Body.String(), "Apikey authorization required\n")
}

func TestAuthRequiredInvalidAuthScheme(t *testing.T) {
	user := newTestUser()

	rw, req, next, rec := mockMiddlewareParams()

	req.Request.Header.Add("Authorization", "Basic YWJjOjEyMw==")

	c, dbs := mockDbContext(&user)
	dbs.Mock.On("GetUser", user.Email).Return(&user, nil)

	ac := new(AuthContext)
	ac.Context = c

	(*AuthContext).AuthRequired(ac, rw, req, next.Next)

	next.Mock.AssertNotCalled(t, "Next", rw, req)
	dbs.Mock.AssertNotCalled(t, "GetUser", user.Email)
	assert.Nil(t, ac.User)
	assert.Equal(t, rec.Code, http.StatusUnauthorized)
	assert.Equal(t, rec.Body.String(), "Apikey authorization required\n")
}

func TestAuthRequiredBadCreds(t *testing.T) {
	user := newTestUser()

	rw, req, next, rec := mockMiddlewareParams()

	req.Request.Header.Add("Authorization", "Apikey asdf")

	c, dbs := mockDbContext(&user)
	dbs.Mock.On("GetUser", user.Email).Return(&user, nil)

	ac := new(AuthContext)
	ac.Context = c

	(*AuthContext).AuthRequired(ac, rw, req, next.Next)

	next.Mock.AssertNotCalled(t, "Next", rw, req)
	dbs.Mock.AssertNotCalled(t, "GetUser", user.Email)
	assert.Nil(t, ac.User)
	assert.Equal(t, rec.Code, http.StatusBadRequest)
	assert.Equal(t, rec.Body.String(), "Invalid authentication params\n")
}

func TestAuthRequiredInvalidUser(t *testing.T) {
	user := newTestUser()

	rw, req, next, rec := mockMiddlewareParams()
	authHeader := fmt.Sprintf("Apikey %s:%s", user.Email, user.ApiKey)
	req.Request.Header.Add("Authorization", authHeader)

	c, dbs := mockDbContext(&user)
	dbs.Mock.On("GetUser", user.Email).Return(nil, errors.New("Invalid user"))

	ac := new(AuthContext)
	ac.Context = c

	(*AuthContext).AuthRequired(ac, rw, req, next.Next)

	next.Mock.AssertNotCalled(t, "Next", rw, req)
	dbs.Mock.AssertCalled(t, "GetUser", user.Email)
	assert.Nil(t, ac.User)
	assert.Equal(t, rec.Code, http.StatusForbidden)
	assert.Equal(t, rec.Body.String(), "Invalid user\n")
}

func TestAuthRequiredInvalidApikey(t *testing.T) {
	user := newTestUser()

	rw, req, next, rec := mockMiddlewareParams()
	authHeader := fmt.Sprintf("Apikey %s:%s", user.Email, user.ApiKey+"!!!")
	req.Request.Header.Add("Authorization", authHeader)

	c, dbs := mockDbContext(&user)
	dbs.Mock.On("GetUser", user.Email).Return(&user, nil)

	ac := new(AuthContext)
	ac.Context = c

	(*AuthContext).AuthRequired(ac, rw, req, next.Next)

	next.Mock.AssertNotCalled(t, "Next", rw, req)
	dbs.Mock.AssertCalled(t, "GetUser", user.Email)
	assert.Nil(t, ac.User)
	assert.Equal(t, rec.Code, http.StatusForbidden)
	assert.Equal(t, rec.Body.String(), "Incorrect API key\n")
}
