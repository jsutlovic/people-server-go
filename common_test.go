package main

import (
	"database/sql"
	"github.com/gocraft/web"
	"github.com/lib/pq/hstore"
	"github.com/stretchr/testify/mock"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
)

func newTestUser() *User {
	user := User{
		1,
		"test@example.com",
		"$2a$04$2a2qnoery/ULUw2WgKVd0OyeHhsHWINab9w9WTPoXqe8xY4PBrwXe",
		"Test User",
		defaultActive,
		defaultSuperuser,
		"",
		nil,
	}

	return &user
}

func newTestPerson(userId int) *Person {
	person := Person{
		1,
		userId,
		"Test Person",
		hstore.Hstore{map[string]sql.NullString{
			"type":  {"asdf", true},
			"other": {"", false},
		}},
		sql.NullInt64{1, true},
		JsonErrors{},
	}

	return &person
}

type mockConfig struct {
	dbType   string
	dbCreds  string
	listener net.Listener
}

func (mc *mockConfig) DbType() string {
	return mc.dbType
}

func (mc *mockConfig) DbCreds() string {
	return mc.dbCreds
}

func (mc *mockConfig) Listener() net.Listener {
	return mc.listener
}

func newTestConfig() Config {
	return &mockConfig{"mock", "", nil}
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

func (m *MockDbService) GetUser(email string) (*User, error) {
	args := m.Mock.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*User), nil
	}
	return nil, args.Error(1)
}

func (m *MockDbService) PasswordCost() int {
	return 4
}

func (m *MockDbService) CreateUser(email, pwhash, name, apikey string, isActive, isSuperuser bool) (*User, error) {
	args := m.Mock.Called(email, pwhash, name, apikey, isActive, isSuperuser)
	if args.Get(0) != nil {
		user := args.Get(0).(*User)
		user.Id = 1
		user.Name = name
		user.Email = email
		user.Pwhash = pwhash
		user.ApiKey = apikey
		user.IsActive = isActive
		user.IsSuperuser = isSuperuser
		return user, nil
	}
	return nil, args.Error(1)
}

func (m *MockDbService) UpdateUser(user *User) error {
	args := m.Mock.Called(user)
	return args.Error(0)
}

func (m *MockDbService) GetPerson(userId, id int) (*Person, error) {
	args := m.Mock.Called(userId, id)
	if args.Get(0) != nil {
		return args.Get(0).(*Person), nil
	}
	return nil, args.Error(1)
}

func (m *MockDbService) GetPeople(userId int) ([]Person, error) {
	args := m.Mock.Called(userId)
	if args.Get(0) != nil {
		return args.Get(0).([]Person), nil
	}
	return nil, args.Error(1)
}

func (m *MockDbService) CreatePerson(userId int, name string, meta hstore.Hstore, color sql.NullInt64) (*Person, error) {
	args := m.Mock.Called(userId, name, meta, color)
	if args.Get(0) != nil {
		person := args.Get(0).(*Person)
		person.Id = 2
		person.UserId = userId
		person.Name = name
		person.Meta = meta
		person.Color = color
		return person, nil
	}
	return nil, args.Error(1)
}

func mockMiddlewareParams() (web.ResponseWriter, *web.Request, *MockNext, *httptest.ResponseRecorder) {
	// Build the ResponseRecorder
	recorder := httptest.NewRecorder()
	rw := new(testResponseWriter)
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

func mockHandlerParams(method, contenttype, content string) (web.ResponseWriter, *web.Request, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	rw := new(testResponseWriter)
	rw.ResponseWriter = recorder

	buf := strings.NewReader(content)
	fakeRequest, err := http.NewRequest(method, "http://example.com/", buf)
	if err != nil {
		panic(err)
	}

	fakeRequest.Header.Set("Content-Type", contenttype)

	req := new(web.Request)
	req.Request = fakeRequest

	return rw, req, recorder
}
