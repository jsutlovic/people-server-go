package main

import (
	"errors"
	"fmt"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
)

const (
	passwordCost     = 10
	defaultActive    = true
	defaultSuperuser = false
)

type UserService interface {
	// User related methods
	GetUser(email string) (*User, error)
	PasswordCost() int
	CreateUser(email, pwhash, name, apikey string, isActive, isSuperuser bool) (*User, error)
}

/*
Basic user struct

Used to model user data to and from the database, as well as to and from the
frontend clients through JSON
*/
type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Pwhash      string `json:"-"`
	Name        string `json:"name"`
	IsActive    bool   `db:"is_active" json:"is_active"`
	IsSuperuser bool   `db:"is_superuser" json:"is_superuser"`
	ApiKey      string `db:"apikey" json:"api_key"`
}

// Compare a given password to this user's current password (hashed)
func (u *User) CheckPassword(password string) bool {
	incorrect := bcrypt.CompareHashAndPassword([]byte(u.Pwhash), []byte(password))

	if incorrect == nil {
		return true
	}
	return false
}

// Compare a given apikey to this user's current apikey
func (u *User) CheckApiKey(apikey string) bool {
	return u.ApiKey == apikey
}

/*
Fetch a user given an email from the database
Returns nil if no matching user is found
*/
func (s *pgDbService) GetUser(email string) (*User, error) {
	user := new(User)

	err := s.db.Get(user, s.db.Rebind(`SELECT * FROM "user" WHERE email=?`), email)
	if err != nil {
		errMsg := fmt.Sprintf("User could not be found: %s", err)
		return nil, errors.New(errMsg)
	}

	return user, nil
}

// Return the default password cost
func (s *pgDbService) PasswordCost() int {
	return passwordCost
}

/*
Create a User in the database with the given email, password hash, name and apikey

IsActive is set to true, IsSuperuser is set to false for the user
*/
func (s *pgDbService) CreateUser(email, pwhash, name, apikey string, isActive, isSuperuser bool) (*User, error) {
	newUser := new(User)

	if strings.TrimSpace(email) == "" {
		return nil, errors.New("Email cannot be empty")
	}

	var userId int

	newUser.Email = email
	newUser.Pwhash = pwhash
	newUser.Name = name
	newUser.IsActive = isActive
	newUser.IsSuperuser = isSuperuser
	newUser.ApiKey = apikey

	insertSql := s.db.Rebind(`INSERT INTO "user" (
        email,
        pwhash,
        name,
        is_active,
        is_superuser,
        apikey
    ) VALUES (?, ?, ?, ?, ?, ?) RETURNING id;`)

	err := s.db.QueryRowx(insertSql,
		newUser.Email,
		newUser.Pwhash,
		newUser.Name,
		newUser.IsActive,
		newUser.IsSuperuser,
		newUser.ApiKey).Scan(&userId)

	if err != nil {
		return nil, err
	}

	newUser.Id = userId

	return newUser, nil
}
