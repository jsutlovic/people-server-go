package main

import (
	"log"

	_ "database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"code.google.com/p/go.crypto/bcrypt"
)

var DB *sqlx.DB

type User struct {
	Id          int
	Email       string
	Pwhash      string
	Name        string
	IsActive    bool `db:"is_active"`
	IsSuperuser bool `db:"is_superuser"`
}

type ApiKey struct {
	UserId int    `db:"user_id"`
	Key    string `db:"key"`
}

func DbInit() {
	DB = sqlx.MustConnect("postgres", "user=vagrant dbname=people host=/var/run/postgresql sslmode=disable application_name=people-go")
	log.Print("db connected")
}

func GetUser(email string) *User {
	var err error
	user := new(User)

	err = DB.Get(user, DB.Rebind("SELECT * FROM \"user\" WHERE email=?"), email)
	if err != nil {
		// We couldn't get the user
		return nil
	}
	log.Printf("Row %d: email %v, password %v\n", user.Id, user.Email, user.Pwhash)
	log.Printf("User struct: %v", user)

	return user
}

func GetApiKey(user User) string {
	var err error
	key := ApiKey{}

	err = DB.Get(&key, DB.Rebind("SELECT * FROM \"api_key\" WHERE user_id=?"), user.Id)
	if err != nil {
		panic(err)
	}

	log.Printf("ApiKey found: %v", key)
	return key.Key
}

func (u *User) CheckPassword(password string) bool {
	incorrect := bcrypt.CompareHashAndPassword([]byte(u.Pwhash), []byte(password))

	if incorrect != nil {
		log.Printf("Password did not match")
		return false
	} else {
		log.Printf("Password matched")
		return true
	}
}
