package main

import (
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
	m.Mock.Called(email)
	return m.User, nil
}

func mockMiddlewareParams() (*web.AppResponseWriter, *web.Request, *MockNext, *httptest.ResponseRecorder) {
	// Build the ResponseRecorder
	recorder := httptest.NewRecorder()
	rw := web.AppResponseWriter{}
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
	next.Mock.On("Next", &rw, &req).Return()

	return &rw, &req, next, recorder
}

func mockDbContext(user *User) (Context, *MockDbService) {
	// Create our mock database service to serve our fake user
	mockDbService := new(MockDbService)
	mockDbService.User = user

	mockDbService.Mock.On("GetUser", user.Email).Return(user)

	// Create Context and set the DB
	c := Context{}
	c.DB = mockDbService

	return c, mockDbService
}

func newTestUser() User {
	// Create a fake user
	user := User{
		1,
		"test@example.com",
		"",
		"Test User",
		true,
		false,
		"abcdefg",
	}

	return user
}

func TestDbMiddleware(t *testing.T) {
	// Get our mock service
	mockDbService := new(MockDbService)

	// Build some basic middleware objects
	rw, req, next, _ := mockMiddlewareParams()

	// Create a Context
	c := Context{}

	// Build the middleware (closure)
	dbMiddleware := DbMiddleware(mockDbService)

	// Call the middleware
	dbMiddleware(&c, rw, req, next.Next)

	// Assertions
	next.Mock.AssertCalled(t, "Next", rw, req)
	assert.Equal(t, c.DB, mockDbService)
}

func TestAuthRequiredAuthorizesValid(t *testing.T) {
	user := newTestUser()

	// Build our basic middleware objects
	rw, req, next, rec := mockMiddlewareParams()

	// Add headers to request
	req.Request.Header.Add("Authorization", fmt.Sprintf("Apikey %s:%s", user.Email, user.ApiKey))

	// Setup contexts
	c, mockDbService := mockDbContext(&user)

	ac := AuthContext{}
	ac.Context = &c

	// Call the middleware
	(*AuthContext).AuthRequired(&ac, rw, req, next.Next)

	// Assertions:

	// Next was called
	next.Mock.AssertCalled(t, "Next", rw, req)
	// GetUser was called
	mockDbService.Mock.AssertCalled(t, "GetUser", user.Email)
	// User is set to the AuthContext
	assert.Equal(t, ac.User, &user)
	// Nothing was written to the responsewriter
	assert.Equal(t, rec.Body.String(), "")
}

func TestAuthRequiredErrorsNoHeader(t *testing.T) {
	user := newTestUser()

	rw, req, next, rec := mockMiddlewareParams()

	c, mockDbService := mockDbContext(&user)

	ac := AuthContext{}
	ac.Context = &c

	(*AuthContext).AuthRequired(&ac, rw, req, next.Next)

	next.Mock.AssertNotCalled(t, "Next", rw, req)
	mockDbService.Mock.AssertNotCalled(t, "GetUser", user.Email)
	assert.Nil(t, ac.User)
	assert.Equal(t, rec.Code, http.StatusUnauthorized)
	assert.Equal(t, rec.Body.String(), "Apikey authorization required\n")
}
