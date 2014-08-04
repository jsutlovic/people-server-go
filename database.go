package main

import (
	"errors"
	"fmt"
	"strings"

	_ "database/sql"
	_ "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"code.google.com/p/go.crypto/bcrypt"
)

const (
	passwordCost = 10
)

/*
Database service

Provides an abstraction wrapper around the database
*/
type DbService interface {
	GetUser(email string) (*User, error)
	PasswordCost() int
	CreateUser(email, pwhash, name, apikey string) (*User, error)
}

type pgDbService struct {
	db *sqlx.DB
}

func NewPgDbService(dbType, creds string) *pgDbService {
	s := new(pgDbService)
	s.dbInit(dbType, creds)
	return s
}

/*
Connect to the database, and set the database handle

This must be called before any database calls can happen

Eventually, this should parse a given config file rather than using
hardcoded values
*/
func (s *pgDbService) dbInit(dbType, creds string) {
	s.db = sqlx.MustConnect(dbType, creds)
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
func (s *pgDbService) CreateUser(email, pwhash, name, apikey string) (*User, error) {
	newUser := new(User)

	if strings.TrimSpace(email) == "" {
		return nil, errors.New("Email cannot be empty")
	}

	var userId int

	newUser.Email = email
	newUser.Pwhash = pwhash
	newUser.Name = name
	newUser.IsActive = true
	newUser.IsSuperuser = false
	newUser.ApiKey = apikey

	insertSql := `INSERT INTO "user" (
		email,
		pwhash,
		name,
		is_active,
		is_superuser,
		apikey
	) VALUES (?, ?, ?, ?, ?, ?) RETURNING id;`

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
