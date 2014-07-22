package main

import (
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
}

func (m *MockDbService) GetUser(email string) (user *User, err error) {
	return new(User), nil
}

func mockMiddlewareParams() (web.ResponseWriter, *web.Request, *MockNext) {
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

	return &rw, &req, next
}

func TestDbMiddleware(t *testing.T) {
	// Get our mock service
	mockDbService := new(MockDbService)

	rw, req, next := mockMiddlewareParams()

	next.Mock.On("Next", rw, req).Return()

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

func TestAuthRequired(t *testing.T) {
}
