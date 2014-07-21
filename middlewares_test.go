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

func TestDbMiddleware(t *testing.T) {
	// Get our mock service
	mockDbService := new(MockDbService)

	// Build the ResponseRecorder
	recorder := httptest.NewRecorder()
	rw := web.AppResponseWriter{}
	rw.ResponseWriter = recorder

	// Build the request
	fakeRequest, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req := web.Request{}
	req.Request = fakeRequest

	// Mock a NextMiddlewareFunc
	mockNext := new(MockNext)

	mockNext.Mock.On("Next", &rw, &req).Return()

	// Create a Context
	c := Context{}

	// Build the middleware (closure)
	dbMiddleware := DbMiddleware(mockDbService)

	// Call the middleware
	dbMiddleware(&c, &rw, &req, mockNext.Next)

	// Assertions
	mockNext.Mock.AssertCalled(t, "Next", &rw, &req)
	assert.Equal(t, c.DB, mockDbService)
}

func TestAuthRequired(t *testing.T) {
}
