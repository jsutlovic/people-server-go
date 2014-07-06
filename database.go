package main

import (
	"log"

	_ "database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"code.google.com/p/go.crypto/bcrypt"
)

var DB *sqlx.DB

func DbInit() {
	DB = sqlx.MustConnect("postgres", "user=vagrant dbname=people host=/var/run/postgresql sslmode=disable application_name=people-go")
	log.Print("db connected")
}

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Pwhash      string `json:"-"`
	Name        string `json:"name"`
	IsActive    bool   `db:"is_active" json:"is_active"`
	IsSuperuser bool   `db:"is_superuser" json:"is_superuser"`
	ApiKey      string `db:"apikey" json:"api_key"`
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
