package main

import (
	"github.com/gocraft/web"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
)

func newTestUser() *User {
	user := User{
		1,
		"test@example.com",
		"",
		"Test User",
		true,
		false,
		"",
	}

	return &user
}

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
	req := new(web.Request)
	req.Request = fakeRequest

	// Mock a NextMiddlewareFunc
	next := new(MockNext)

	// Setup expecations for Next
	next.Mock.On("Next", rw, req).Return()

	return rw, req, next, recorder
}

func mockDbContext(user *User) (*Context, *MockDbService) {
	// Create our mock database service to serve our fake user
	dbs := new(MockDbService)
	dbs.User = user

	// Create Context and set the DB
	c := new(Context)
	c.DB = dbs

	return c, dbs
}

func mockAuthContext(user *User) (*AuthContext, *MockDbService) {
	c, dbs := mockDbContext(user)

	ac := new(AuthContext)
	ac.Context = c
	ac.User = user

	return ac, dbs
}

func mockHandlerParams(method string) (web.ResponseWriter, *web.Request, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	rw := new(web.AppResponseWriter)
	rw.ResponseWriter = recorder

	fakeRequest, err := http.NewRequest(method, "http://example.com/", nil)
	if err != nil {
		panic(err)
	}

	req := new(web.Request)
	req.Request = fakeRequest

	return rw, req, recorder
}
