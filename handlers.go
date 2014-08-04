package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocraft/web"
	"net/http"
	"regexp"
	"strings"
)

const (
	JsonContentType = "application/json"

	// Auth errors
	InvalidCredentials = "Invalid credentials"
	ParamsRequired     = "Email and password are required"
	InactiveUser       = "User disabled"

	// UserCreate.Validate errors
	UserCreateEmailEmpty     = "Email cannot be empty"
	UserCreatePasswordEmpty  = "Password cannot be empty"
	UserCreateNameEmpty      = "Name cannot be empty"
	UserCreatePasswordLength = "Password must be at least 6 characters long"
	UserCreateInvalidEmail   = "Invalid email address"
	UserCreateError          = "Error creating user"

	// UserCreateApi errors
	JsonContentTypeError = "Content-Type is not JSON"
	JsonMalformedError   = "Malformed JSON"
	InvalidUserDataError = "Invalid User data"
	UserExistsError      = "User already exists"
)

// JSON structs for web APIs

type UserErrors map[string]string

// Expected format of JSON data for creating a User
type UserCreate struct {
	Email    string     `json:"email"`
	Password string     `json:"password"`
	Name     string     `json:"name"`
	errors   UserErrors `json:"-"`
}

func (u *UserCreate) Errors() UserErrors {
	if u.errors == nil {
		u.errors = UserErrors{}
	}
	return u.errors
}

func (u *UserCreate) Validate() bool {
	anyBlank := false
	fieldErrors := false

	u.errors = UserErrors{}

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
			rw.Header().Set("Content-Type", JsonContentType)
			fmt.Fprint(rw, Jsonify(user))
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
	rw.Header().Set("Content-Type", JsonContentType)
	fmt.Fprint(rw, Jsonify(c.User))
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

	rw.Header().Set("Content-Type", JsonContentType)
	fmt.Fprint(rw, Jsonify(user))
}
