package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	JsonContentType = "application/json"

	// Auth errors
	InvalidCredentials = "Invalid credentials"
	InvalidUser        = "Invalid user"
	ParamsRequired     = "Email and password are required"
	InactiveUser       = "User disabled"

	// UserCreate.Validate errors
	UserCreateEmailEmpty     = "Email cannot be empty"
	UserCreatePasswordEmpty  = "Password cannot be empty"
	UserCreateNameEmpty      = "Name cannot be empty"
	UserCreatePasswordLength = "Password must be at least 6 characters long"
	UserCreateInvalidEmail   = "Invalid email address"

	// UserCreateApi errors
	UserCreateError      = "Error creating user"
	JsonContentTypeError = "Content-Type is not JSON"
	JsonMalformedError   = "Malformed JSON"
	InvalidUserDataError = "Invalid User data"
	UserExistsError      = "User already exists"

	// PersonCreateApi errors
	PersonCreateError = "Error creating person"
)

// Basic Context available to all handlers
type Context struct {
	DB DbService
}

// Context supplying an authorized user. Used with AuthRequired middleware
type AuthContext struct {
	*Context
	User *User
}

// JSON structs for web APIs

type JsonErrors map[string]string

// Expected format of JSON data for creating a User
type UserCreate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	errors   JsonErrors
}

func (u *UserCreate) Errors() JsonErrors {
	if u.errors == nil {
		u.errors = JsonErrors{}
	}
	return u.errors
}

// Validation for UserCreate struct
func (u *UserCreate) Validate() bool {
	anyBlank := false
	fieldErrors := false

	u.errors = JsonErrors{}

	email := strings.TrimSpace(u.Email)
	name := strings.TrimSpace(u.Name)

	if email == "" {
		u.errors["email"] = UserCreateEmailEmpty
		anyBlank = true
	}
	// Spaces are valid in passwords
	if u.Password == "" {
		u.errors["password"] = UserCreatePasswordEmpty
		anyBlank = true
	}
	if name == "" {
		u.errors["name"] = UserCreateNameEmpty
		anyBlank = true
	}
	if anyBlank {
		return false
	}

	if len(u.Password) < 4 {
		u.errors["password"] = UserCreatePasswordLength
		fieldErrors = true
	}

	re := regexp.MustCompile(".+@.+\\..+")
	matched := re.Match([]byte(email))
	if !matched {
		u.errors["email"] = UserCreateInvalidEmail
		fieldErrors = true
	}

	return !fieldErrors
}

// Expected format of JSON data for a Person
type PersonJSON struct {
	Id     int               `json:"id"`
	UserId int               `json:"user_id"`
	Name   string            `json:"name"`
	Meta   map[string]string `json:"meta"`
	Color  json.RawMessage   `json:"color"`
}

func jsonResponse(rw web.ResponseWriter, data interface{}) {
	rw.Header().Set("Content-Type", JsonContentType)
	fmt.Fprint(rw, Jsonify(data))
}

/*
Handler to authenticate a user.

Parses the current request form looking for 'email' and 'password'
fields. Checks the database for the given user and password.

If authentication is successful, return a JSON object representing the
authenticated user. Otherwise returns 403 Forbidden.
*/
func (c *Context) ApiAuth(rw web.ResponseWriter, req *web.Request) {
	req.ParseForm()

	form := req.PostForm

	emails, email_ok := form["email"]
	passwords, password_ok := form["password"]
	if !(email_ok && password_ok) {
		status := http.StatusBadRequest
		rw.WriteHeader(status)
		fmt.Fprint(rw, ParamsRequired)
		return
	}

	email := emails[0]
	password := passwords[0]

	user, err := c.DB.GetUser(email)
	authed := err == nil && user != nil && user.CheckPassword(password)
	if authed {
		if user.IsActive {
			jsonResponse(rw, user)
		} else {
			http.Error(rw, "User disabled", http.StatusForbidden)
		}
	} else {
		http.Error(rw, InvalidCredentials, http.StatusForbidden)
	}
}

/*
Handler for the GET User API

Returns a JSON representation of the currently authenticated User
*/
func (c *AuthContext) GetUserApi(rw web.ResponseWriter, req *web.Request) {
	jsonResponse(rw, c.User)
}

/*
Handler for the POST User API

Takes a JSON representation of a user and creates an account for it
If a duplicate email is found, return 409 Conflict
Otherwise, if creation is successful return 201 Created
*/
func (c *Context) CreateUserApi(rw web.ResponseWriter, req *web.Request) {
	ct, ctok := req.Header["Content-Type"]
	if !ctok || len(ct) < 1 || (len(ct) >= 1 && ct[0] != "application/json") {
		http.Error(rw, JsonContentTypeError, http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(req.Body)
	newUser := new(UserCreate)
	err := dec.Decode(&newUser)
	if err != nil {
		http.Error(rw, JsonMalformedError, http.StatusBadRequest)
		return
	}

	if !newUser.Validate() {
		http.Error(rw, Jsonify(newUser.Errors()), http.StatusBadRequest)
		return
	}

	if existing, _ := c.DB.GetUser(newUser.Email); existing != nil {
		http.Error(rw, UserExistsError, http.StatusConflict)
		return
	}

	user, err := c.DB.CreateUser(
		newUser.Email,
		GeneratePasswordHash(newUser.Password, c.DB.PasswordCost()),
		newUser.Name,
		GenerateApiKey())

	if err != nil {
		http.Error(rw, UserCreateError, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	jsonResponse(rw, user)
}

/*
Handler for GET Person API

Returns a single Person as JSON, or 404 if the user does not have access to that Person
*/
func (c *AuthContext) GetPersonApi(rw web.ResponseWriter, req *web.Request) {
	idStr, idExists := req.PathParams["id"]

	if !idExists {
		http.Error(rw, "Invalid Path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(rw, "id must be an integer", http.StatusBadRequest)
		return
	}

	person, err := c.DB.GetPerson(c.User.Id, id)

	if err != nil {
		http.Error(rw, "Person not found", http.StatusNotFound)
		return
	}

	jsonResponse(rw, person)
}

/*
Handler for GET Person List API

Returns all the Person objects associated with the current User
*/
func (c *AuthContext) GetPersonListApi(rw web.ResponseWriter, req *web.Request) {
	people, err := c.DB.GetPeople(c.User.Id)

	if err != nil {
		http.Error(rw, "No people found", http.StatusNotFound)
		return
	}

	jsonResponse(rw, people)
}

/*
Handler for POST Person API

Takes a JSON representation of a Person and creates it for the user logged in.
*/
func (c *AuthContext) CreatePersonApi(rw web.ResponseWriter, req *web.Request) {
	if c.User == nil || c.User.Id == 0 {
		http.Error(rw, InvalidUser, http.StatusUnauthorized)
		return
	}

	ct, ctok := req.Header["Content-Type"]
	if !ctok || len(ct) < 1 || (len(ct) >= 1 && ct[0] != "application/json") {
		http.Error(rw, JsonContentTypeError, http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(req.Body)
	newPerson := new(Person)
	err := dec.Decode(&newPerson)
	if err != nil {
		http.Error(rw, JsonMalformedError, http.StatusBadRequest)
		return
	}

	if !newPerson.Validate() {
		http.Error(rw, Jsonify(newPerson.Errors()), http.StatusBadRequest)
		return
	}
}
