package main

import (
	"errors"
	"log"

	_ "database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"code.google.com/p/go.crypto/bcrypt"
)

/*
Database service

Provides an abstraction wrapper around the database
*/
type DbService interface {
	GetUser(email string) (user *User, err error)
}

type pgDbService struct {
	db *sqlx.DB
}

func NewPgDbService() *pgDbService {
	s := new(pgDbService)
	s.dbInit()
	return s
}

/*
Connect to the database, and set the database handle

This must be called before any database calls can happen

Eventually, this should parse a given config file rather than using
hardcoded values
*/
func (s *pgDbService) dbInit() {
	s.db = sqlx.MustConnect("postgres", "user=vagrant dbname=people host=/var/run/postgresql sslmode=disable application_name=people-go")
	log.Print("db connected")
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
func (s *pgDbService) GetUser(email string) (user *User, err error) {
	user = new(User)

	err = s.db.Get(user, s.db.Rebind("SELECT * FROM \"user\" WHERE email=?"), email)
	if err != nil {
		// We couldn't get the user
		log.Println(err)
		return nil, errors.New("User could not be found")
	}

	return user, nil
}

// Compare a given password to this user's current password (hashed)
func (u *User) CheckPassword(password string) bool {
	incorrect := bcrypt.CompareHashAndPassword([]byte(u.Pwhash), []byte(password))

	if incorrect == nil {
		log.Printf("Password matched")
		return true
	}
	log.Printf("Password did not match")
	return false
}

// Compare a given apikey to this user's current apikey
func (u *User) CheckApiKey(apikey string) bool {
	return u.ApiKey == apikey
}
