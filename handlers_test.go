package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
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

	dbs.Mock.AssertExpectations(t)
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

	dbs.Mock.AssertExpectations(t)
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

	dbs.Mock.AssertExpectations(t)
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

	dbs.Mock.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusForbidden)
	assert.Equal(t, rec.Body.String(), InactiveUser+"\n")
}

func TestUserCreateFields(t *testing.T) {
	uc := UserCreate{}
	userCreateType := reflect.TypeOf(uc)

	fieldCount := userCreateType.NumField()
	assert.Equal(t, fieldCount, 4)

	_, emailExists := userCreateType.FieldByName("Email")
	_, passwordExists := userCreateType.FieldByName("Password")
	_, nameExists := userCreateType.FieldByName("Name")
	_, errorExists := userCreateType.FieldByName("errors")

	assert.True(t, emailExists)
	assert.True(t, passwordExists)
	assert.True(t, nameExists)
	assert.True(t, errorExists)
}

func TestUserCreateJsonTags(t *testing.T) {
	uc := UserCreate{}
	userCreateType := reflect.TypeOf(uc)

	emailField, _ := userCreateType.FieldByName("Email")
	passwordField, _ := userCreateType.FieldByName("Password")
	nameField, _ := userCreateType.FieldByName("Name")
	errorsField, _ := userCreateType.FieldByName("errors")

	assert.Equal(t, emailField.Tag.Get("json"), "email")
	assert.Equal(t, passwordField.Tag.Get("json"), "password")
	assert.Equal(t, nameField.Tag.Get("json"), "name")
	assert.Equal(t, errorsField.Tag.Get("json"), "-")
}

func TestUserCreateValidateCreatesNewErrors(t *testing.T) {
	uc := UserCreate{}

	nilErrors := uc.Errors()
	assert.NotNil(t, nilErrors, "Calling Errors() does not fail if errors is not initialized")

	uc.Validate()

	firstErrors := uc.Errors()
	secondErrors := uc.Errors()

	assert.Equal(t, firstErrors, secondErrors)

	uc.Email = "test@example.com"
	uc.Validate()

	thirdErrors := uc.Errors()
	assert.NotEqual(t, firstErrors, thirdErrors)
}

func TestUserCreateValidate(t *testing.T) {
	validateTests := []struct {
		in     UserCreate
		out    bool
		errors UserErrors
	}{
		// Invalid
		{
			in:  UserCreate{"", "", "", nil},
			out: false,
			errors: UserErrors{
				"email":    UserCreateEmailEmpty,
				"password": UserCreatePasswordEmpty,
				"name":     UserCreateNameEmpty,
			},
		},
		{
			in:  UserCreate{" ", " ", " ", nil},
			out: false,
			errors: UserErrors{
				"email": UserCreateEmailEmpty,
				"name":  UserCreateNameEmpty,
			},
		},
		{
			in:  UserCreate{"\t", "", "\t ", nil},
			out: false,
			errors: UserErrors{
				"email":    UserCreateEmailEmpty,
				"password": UserCreatePasswordEmpty,
				"name":     UserCreateNameEmpty,
			},
		},
		{
			in:  UserCreate{" \n\r ", "", "\r\n ", nil},
			out: false,
			errors: UserErrors{
				"email":    UserCreateEmailEmpty,
				"password": UserCreatePasswordEmpty,
				"name":     UserCreateNameEmpty,
			},
		},
		{
			in:  UserCreate{" \n\t\r ", "", "\r\n ", nil},
			out: false,
			errors: UserErrors{
				"email":    UserCreateEmailEmpty,
				"password": UserCreatePasswordEmpty,
				"name":     UserCreateNameEmpty,
			},
		},
		{
			in:  UserCreate{"test@example.com", "", "", nil},
			out: false,
			errors: UserErrors{
				"password": UserCreatePasswordEmpty,
				"name":     UserCreateNameEmpty,
			},
		},
		{
			in:  UserCreate{"", "asdf", "", nil},
			out: false,
			errors: UserErrors{
				"email": UserCreateEmailEmpty,
				"name":  UserCreateNameEmpty,
			},
		},
		{
			in:  UserCreate{"", "", "Test User", nil},
			out: false,
			errors: UserErrors{
				"email":    UserCreateEmailEmpty,
				"password": UserCreatePasswordEmpty,
			},
		},
		{
			in:  UserCreate{"test@example.com", "asdf", "", nil},
			out: false,
			errors: UserErrors{
				"name": UserCreateNameEmpty,
			},
		},
		{
			in:  UserCreate{"", "asdf", "Test User", nil},
			out: false,
			errors: UserErrors{
				"email": UserCreateEmailEmpty,
			},
		},
		{
			in:  UserCreate{"test@example.com", "", "Test User", nil},
			out: false,
			errors: UserErrors{
				"password": UserCreatePasswordEmpty,
			},
		},
		{
			in:  UserCreate{"test@example.com", "asd", "Test User", nil},
			out: false,
			errors: UserErrors{
				"password": UserCreatePasswordLength,
			},
		},
		{
			in:  UserCreate{"test", "asdf", "Test User", nil},
			out: false,
			errors: UserErrors{
				"email": UserCreateInvalidEmail,
			},
		},
		{
			in:  UserCreate{"test@example", "asdf", "Test User", nil},
			out: false,
			errors: UserErrors{
				"email": UserCreateInvalidEmail,
			},
		},
		{
			in:  UserCreate{"@example", "asdf", "Test User", nil},
			out: false,
			errors: UserErrors{
				"email": UserCreateInvalidEmail,
			},
		},
		{
			in:  UserCreate{"@example.com", "asdf", "Test User", nil},
			out: false,
			errors: UserErrors{
				"email": UserCreateInvalidEmail,
			},
		},
		{
			in:  UserCreate{"@example.co.uk", "asdf", "Test User", nil},
			out: false,
			errors: UserErrors{
				"email": UserCreateInvalidEmail,
			},
		},
		{
			in:  UserCreate{"example.com", "asdf", "Test User", nil},
			out: false,
			errors: UserErrors{
				"email": UserCreateInvalidEmail,
			},
		},
		{
			in:  UserCreate{"test@example", "a", "Test User", nil},
			out: false,
			errors: UserErrors{
				"email":    UserCreateInvalidEmail,
				"password": UserCreatePasswordLength,
			},
		},
		{
			in:  UserCreate{"test@example", "a", "", nil},
			out: false,
			errors: UserErrors{
				"name": UserCreateNameEmpty,
			},
		},
		{
			in:  UserCreate{"", "a", "Test User", nil},
			out: false,
			errors: UserErrors{
				"email": UserCreateEmailEmpty,
			},
		},
		{
			in:  UserCreate{"test@example", "", "Test User", nil},
			out: false,
			errors: UserErrors{
				"password": UserCreatePasswordEmpty,
			},
		},

		// Valid
		{
			in:     UserCreate{"test@example.com", "asdf", "Test User", nil},
			out:    true,
			errors: UserErrors{},
		},
	}

	for _, test := range validateTests {
		actualValid := test.in.Validate()

		assert.Equal(t, actualValid, test.out)
		assert.Equal(t, test.in.Errors(), test.errors)
	}
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
	badformats := []struct {
		in  map[string]string
		out map[string]string
	}{
		{
			in: map[string]string{
				"user": "test@example.com",
			},
			out: map[string]string{
				"email":    UserCreateEmailEmpty,
				"password": UserCreatePasswordEmpty,
				"name":     UserCreateNameEmpty,
			},
		},
		{
			in: map[string]string{
				"email": "test@example.com",
			},
			out: map[string]string{
				"password": UserCreatePasswordEmpty,
				"name":     UserCreateNameEmpty,
			},
		},
		{
			in: map[string]string{
				"email": "test@example.com",
				"name":  "Test User",
			},
			out: map[string]string{
				"password": UserCreatePasswordEmpty,
			},
		},
		{
			in: map[string]string{
				"email":    "test@example.com",
				"password": "asdf",
			},
			out: map[string]string{
				"name": UserCreateNameEmpty,
			},
		},
		{
			in: map[string]string{
				"password": "asdf",
				"name":     "Test User",
			},
			out: map[string]string{
				"email": UserCreateEmailEmpty,
			},
		},
		{
			in: map[string]string{},
			out: map[string]string{
				"email":    UserCreateEmailEmpty,
				"password": UserCreatePasswordEmpty,
				"name":     UserCreateNameEmpty,
			},
		},
		{
			in: map[string]string{
				"username": "testuser",
			},
			out: map[string]string{
				"email":    UserCreateEmailEmpty,
				"password": UserCreatePasswordEmpty,
				"name":     UserCreateNameEmpty,
			},
		},
		{
			in: map[string]string{
				"email":    "test",
				"password": "asdf",
				"name":     "Test User",
			},
			out: map[string]string{
				"email": UserCreateInvalidEmail,
			},
		},
		{
			in: map[string]string{
				"email":    "test@example.com",
				"password": "a",
				"name":     "Test User",
			},
			out: map[string]string{
				"password": UserCreatePasswordLength,
			},
		},
		{
			in: map[string]string{
				"email":    "test@example",
				"password": "a",
				"name":     "Test User",
			},
			out: map[string]string{
				"email":    UserCreateInvalidEmail,
				"password": UserCreatePasswordLength,
			},
		},
	}

	for _, test := range badformats {
		rw, req, rec := mockHandlerParams("POST", JsonContentType, Jsonify(test.in))

		c, _ := mockDbContext(nil)

		(*Context).CreateUserApi(c, rw, req)

		assert.Equal(t, rec.Code, http.StatusBadRequest)

		dec := json.NewDecoder(rec.Body)
		var actualOutJson map[string]string
		dec.Decode(&actualOutJson)
		assert.Equal(t, actualOutJson, test.out)
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
		rw, req, rec := mockHandlerParams("POST", JsonContentType, data)

		c, _ := mockDbContext(nil)

		(*Context).CreateUserApi(c, rw, req)

		assert.Equal(t, rec.Code, http.StatusBadRequest)
		assert.Equal(t, rec.Body.String(), JsonMalformedError+"\n")
	}
}

func TestCreateUserApiUserExists(t *testing.T) {
	newUser := UserCreate{"test@example.com", "asdf", "Test User", nil}
	user := newTestUser()
	user.Email = newUser.Email

	rw, req, rec := mockHandlerParams("POST", JsonContentType, Jsonify(newUser))

	c, dbs := mockDbContext(user)

	dbs.Mock.On("GetUser", user.Email).Return(user, nil)

	(*Context).CreateUserApi(c, rw, req)

	assert.Equal(t, rec.Code, http.StatusConflict)
	assert.Equal(t, rec.Body.String(), UserExistsError+"\n")
}

func TestCreateUserApiInsertError(t *testing.T) {
	userEmail := "test@example.com"
	userPassword := "asdf"
	userName := "Test User"
	newUser := UserCreate{userEmail, userPassword, userName, nil}

	rw, req, rec := mockHandlerParams("POST", JsonContentType, Jsonify(newUser))

	c, dbs := mockDbContext(nil)

	dbs.Mock.On("GetUser", userEmail).Return(nil, errors.New("No user found"))

	pwhash := mock.AnythingOfType("string")
	apikey := mock.AnythingOfType("string")
	dbs.Mock.On("CreateUser", userEmail, pwhash, userName, apikey).Return(nil, errors.New(""))

	(*Context).CreateUserApi(c, rw, req)

	assert.Equal(t, rec.Code, http.StatusInternalServerError)
	assert.Equal(t, rec.Body.String(), UserCreateError+"\n")
}

func TestCreateUserApi(t *testing.T) {
	userEmail := "test@example.com"
	userPassword := "asdf"
	userName := "Test User"
	newUser := UserCreate{userEmail, userPassword, userName, nil}
	user := newTestUser()
	user.Email = userEmail
	user.Name = userName

	rw, req, rec := mockHandlerParams("POST", JsonContentType, Jsonify(newUser))

	c, dbs := mockDbContext(nil)

	dbs.Mock.On("GetUser", userEmail).Return(nil, errors.New("No user found"))

	pwhash := mock.AnythingOfType("string")
	apikey := mock.AnythingOfType("string")
	dbs.Mock.On("CreateUser", userEmail, pwhash, userName, apikey).Return(user, nil)

	(*Context).CreateUserApi(c, rw, req)

	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), Jsonify(user))
	assert.True(t, user.CheckPassword(userPassword))
	assert.Len(t, user.ApiKey, 40)
}

func TestGetPersonApiNoId(t *testing.T) {
	userId := 2
	personId := 1

	rw, req, rec := mockHandlerParams("GET", "", "")

	user := newTestUser()
	user.Id = userId

	ac, dbs := mockAuthContext(user)

	(*AuthContext).GetPersonApi(ac, rw, req)

	dbs.Mock.AssertNotCalled(t, "GetPerson", userId, personId)
	assert.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestGetPersonApiInvalidId(t *testing.T) {
	userId := 2
	personId := 1

	rw, req, rec := mockHandlerParams("GET", "", "")
	req.PathParams = map[string]string{"id": "sadf"}

	user := newTestUser()
	user.Id = userId

	ac, dbs := mockAuthContext(user)

	(*AuthContext).GetPersonApi(ac, rw, req)

	dbs.Mock.AssertNotCalled(t, "GetPerson", userId, personId)
	assert.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestGetPersonApiNonExisting(t *testing.T) {
	userId := 2
	personId := 1

	rw, req, rec := mockHandlerParams("GET", "", "")
	req.PathParams = map[string]string{"id": strconv.Itoa(personId)}

	user := newTestUser()
	user.Id = userId

	ac, dbs := mockAuthContext(user)

	dbs.Mock.On("GetPerson", userId, personId).Return(nil, errors.New("Not found"))

	(*AuthContext).GetPersonApi(ac, rw, req)

	dbs.Mock.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusNotFound)
}

func TestGetPersonApi(t *testing.T) {
	userId := 2
	personId := 1

	rw, req, rec := mockHandlerParams("GET", "", "")
	req.PathParams = map[string]string{"id": strconv.Itoa(personId)}

	user := newTestUser()
	user.Id = userId
	person := newTestPerson(userId)

	ac, dbs := mockAuthContext(user)

	dbs.Mock.On("GetPerson", userId, personId).Return(person, nil)

	(*AuthContext).GetPersonApi(ac, rw, req)

	dbs.Mock.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), Jsonify(person))
}

func TestGetPersonListApiError(t *testing.T) {
	userId := 1

	rw, req, rec := mockHandlerParams("GET", "", "")

	user := newTestUser()
	user.Id = userId

	ac, dbs := mockAuthContext(user)

	dbs.Mock.On("GetPeople", userId).Return(nil, errors.New("DB Error"))

	(*AuthContext).GetPersonListApi(ac, rw, req)

	dbs.Mock.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusNotFound)
}

func TestGetPersonListApiEmpty(t *testing.T) {
	userId := 1

	rw, req, rec := mockHandlerParams("GET", "", "")

	user := newTestUser()
	user.Id = userId

	ac, dbs := mockAuthContext(user)

	pp := []Person{}

	dbs.Mock.On("GetPeople", userId).Return(pp, nil)

	(*AuthContext).GetPersonListApi(ac, rw, req)

	dbs.Mock.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), "[]")
}

func TestGetPersonListApi(t *testing.T) {
	userId := 1

	rw, req, rec := mockHandlerParams("GET", "", "")

	user := newTestUser()
	user.Id = userId

	p1 := newTestPerson(userId)
	p1.Id = 1
	p1.Name = "Test 1"

	p2 := newTestPerson(userId)
	p2.Id = 2
	p2.Name = "Test 2"

	ac, dbs := mockAuthContext(user)

	pp := []Person{*p1, *p2}

	dbs.Mock.On("GetPeople", userId).Return(pp, nil)

	(*AuthContext).GetPersonListApi(ac, rw, req)

	dbs.Mock.AssertExpectations(t)
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), Jsonify(pp))
}
