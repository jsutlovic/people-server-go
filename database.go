package main

import (
	_ "database/sql"
	_ "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

/*
Database service

Provides an abstraction wrapper around the database
*/
type DbService interface {
	// User related methods
	GetUser(email string) (*User, error)
	PasswordCost() int
	CreateUser(email, pwhash, name, apikey string) (*User, error)

	// People related methods
	GetPerson(userId, id int) (*Person, error)
	GetPeople(userId int) ([]Person, error)
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
