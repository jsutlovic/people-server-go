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

func CheckPassword(email string, password string) bool {
	var err error
	var ul []User

	err = DB.Select(&ul, DB.Rebind("SELECT id, email, pwhash, name, is_active, is_superuser FROM \"user\" WHERE email=?"), email)
	if err != nil {
		panic(err)
	}
	if len(ul) < 1 {
		log.Println("User not found")
	}
	for _, user := range ul {
		log.Printf("Row %d: email %v, password %v\n", user.Id, user.Email, user.Pwhash)
		log.Printf("User struct: %v", user)

		incorrect := bcrypt.CompareHashAndPassword([]byte(user.Pwhash), []byte(password))

		if incorrect != nil {
			log.Printf("Password did not match")
		} else {
			log.Printf("Password matched")
			return true
		}
	}
	return false
}
