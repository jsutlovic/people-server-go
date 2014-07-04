package main

type User struct {
	Id          int
	Email       string
	Pwhash      string
	Name        string
	IsActive    bool `db:"is_active"`
	IsSuperuser bool `db:"is_superuser"`
}
