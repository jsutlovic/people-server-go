package main

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestUserCheckPassword(t *testing.T) {
	var password string = "asdf"
	var hashedInputs []string = []string{
		"$2a$04$2a2qnoery/ULUw2WgKVd0OyeHhsHWINab9w9WTPoXqe8xY4PBrwXe",
		"$2a$05$grq2.qlk2BmVFHQ8Uih/L.qt7rJjJjpgsEmVw7BBeIqjCid9UTCxe",
		"$2a$06$7A7qeAPNl/4jcYvzsYRngudI.MdeHh944QU/25fcecrIxnyTv4ria",
		"$2a$07$CBk.ZLMUbZczyH1MOObo1em9kMaue/MFdIg.0vNBzefCInMbZ.hRK",
		"$2a$08$6FSu6citVZco8DldlFztoegK1Q0LWQ66Nu5MlHUb.R6ocj2UEp3Cy",
		"$2a$09$Bnm489iWApygcN2SObO5u.U.HGnOSXi5UuoNfEc3eyorLf218KVnu",
		"$2a$10$Favll94j6lFbNU4iLgFlDe.PRZNzNBmK.I7vU15bmulv2RCLAFpRK",
	}
	var badInputs []string = []string{
		"asdf",
		"pbkdf2_sha256$12000$8POIxt1QfIjM$Jrsr61tAHdITnf7NhiXg/MaSHn0k/sczKOZGdnEmPFc=",
		"$2a$04$lVOHxIVERv4/e3Ch0opdFeBNQC1mX4FMrgiHY4EjvDS4EhKwmqQsO",
		"$2a$05$dQa8ul8msIij9WF6YhNl0.nGp4PBdNZ36EzOW7YmkiwnBJEe.eUcG",
		"$2a$04$jyUNJGCcDdKEr/K.9JwK.u9jFcMFc8ZQ9j2sQLuB5Ge4QxNTKRbBS",
	}

	user := newTestUser()

	for _, hashedInput := range hashedInputs {
		user.Pwhash = hashedInput
		assert.True(t, user.CheckPassword(password))
	}

	for _, badInput := range badInputs {
		user.Pwhash = badInput
		assert.False(t, user.CheckPassword(password))
	}
}

func TestUserCheckApiKey(t *testing.T) {
	var apikeyInputs []string = []string{
		"7fc61f88d375a2d784e20034a82cb95a7e4a589c",
		"115ba60d24ea754a7f1f940680f18669b36a717f",
		"d366ef943b2ff274b39ed330703adb13be32a5a5",
	}

	user := newTestUser()

	lastInput := ""
	for _, apikey := range apikeyInputs {
		user.ApiKey = apikey

		assert.True(t, user.CheckApiKey(apikey))
		assert.False(t, user.CheckApiKey(lastInput))
		lastInput = apikey
	}
}

func TestUserFields(t *testing.T) {
	u := User{}
	userType := reflect.TypeOf(u)

	fieldCount := userType.NumField()
	assert.Equal(t, fieldCount, 7)
	_, idExists := userType.FieldByName("Id")
	_, emailExists := userType.FieldByName("Email")
	_, pwExists := userType.FieldByName("Pwhash")
	_, nameExists := userType.FieldByName("Name")
	_, activeExists := userType.FieldByName("IsActive")
	_, superExists := userType.FieldByName("IsSuperuser")
	_, apikeyExists := userType.FieldByName("ApiKey")

	assert.True(t, idExists)
	assert.True(t, emailExists)
	assert.True(t, pwExists)
	assert.True(t, nameExists)
	assert.True(t, activeExists)
	assert.True(t, activeExists)
	assert.True(t, superExists)
	assert.True(t, apikeyExists)
}

func TestUserJsonTags(t *testing.T) {
	u := User{}
	userType := reflect.TypeOf(u)

	idField, _ := userType.FieldByName("Id")
	emailField, _ := userType.FieldByName("Email")
	pwField, _ := userType.FieldByName("Pwhash")
	nameField, _ := userType.FieldByName("Name")
	activeField, _ := userType.FieldByName("IsActive")
	superField, _ := userType.FieldByName("IsSuperuser")
	apikeyField, _ := userType.FieldByName("ApiKey")

	assert.Equal(t, idField.Tag.Get("json"), "id")
	assert.Equal(t, emailField.Tag.Get("json"), "email")
	assert.Equal(t, pwField.Tag.Get("json"), "-")
	assert.Equal(t, nameField.Tag.Get("json"), "name")
	assert.Equal(t, activeField.Tag.Get("json"), "is_active")
	assert.Equal(t, superField.Tag.Get("json"), "is_superuser")
	assert.Equal(t, apikeyField.Tag.Get("json"), "api_key")
}

func TestUserDbTags(t *testing.T) {
	u := User{}
	userType := reflect.TypeOf(u)

	idField, _ := userType.FieldByName("Id")
	emailField, _ := userType.FieldByName("Email")
	pwField, _ := userType.FieldByName("Pwhash")
	nameField, _ := userType.FieldByName("Name")
	activeField, _ := userType.FieldByName("IsActive")
	superField, _ := userType.FieldByName("IsSuperuser")
	apikeyField, _ := userType.FieldByName("ApiKey")

	assert.Equal(t, idField.Tag.Get("db"), "")
	assert.Equal(t, emailField.Tag.Get("db"), "")
	assert.Equal(t, pwField.Tag.Get("db"), "")
	assert.Equal(t, nameField.Tag.Get("db"), "")
	assert.Equal(t, activeField.Tag.Get("db"), "is_active")
	assert.Equal(t, superField.Tag.Get("db"), "is_superuser")
	assert.Equal(t, apikeyField.Tag.Get("db"), "apikey")
}

func TestGetUser(t *testing.T) {
	pgdbs := NewPgDbService("mock", "")

	userEmail := "test@example.com"
	cols := []string{"id", "email", "pwhash", "name", "is_active", "is_superuser", "apikey"}
	data := "1,test@example.com,,Test User,true,false,abcdefg"

	sqlmock.ExpectQuery(`SELECT \* FROM "user" WHERE email=?`).
		WithArgs(userEmail).
		WillReturnRows(sqlmock.NewRows(cols).FromCSVString(data))

	u, err := pgdbs.GetUser(userEmail)
	if !assert.Nil(t, err, "Query should not error") {
		return
	}

	if !assert.NotNil(t, u, "User should not be nil") {
		return
	}

	assert.Equal(t, u.Email, userEmail)
}

func TestGetUserError(t *testing.T) {
	pgdbs := NewPgDbService("mock", "")

	userEmail := "test2@example.com"

	sqlmock.ExpectQuery(`SELECT \* FROM "user" WHERE email=?`).
		WithArgs(userEmail).
		WillReturnError(errors.New("Could not find user"))

	u, err := pgdbs.GetUser(userEmail)
	if !assert.Nil(t, u, "User should be nil") {
		return
	}

	if !assert.NotNil(t, err, "Should have an error") {
		return
	}

	assert.Equal(t, err.Error(), "User could not be found: Could not find user")
}

func TestCreateUserInsertError(t *testing.T) {
	pgdbs := NewPgDbService("mock", "")

	userEmail := "test@example.com"
	userPwhash := "$2a$04$2a2qnoery/ULUw2WgKVd0OyeHhsHWINab9w9WTPoXqe8xY4PBrwXe"
	userName := "Test User"
	userApikey := GenerateApiKey()

	sqlmock.ExpectQuery(`INSERT INTO "user" \( email, pwhash, name, is_active, is_superuser, apikey \) VALUES \(\?, \?, \?, \?, \?, \?\) RETURNING id;`).
		WithArgs(userEmail, userPwhash, userName, true, false, userApikey).
		WillReturnError(errors.New("Could not insert"))

	u, err := pgdbs.CreateUser(userEmail, userPwhash, userName, userApikey)

	if !assert.Nil(t, u, "User should be nil") {
		return
	}

	if !assert.NotNil(t, err, "Error should not be nil") {
		return
	}

	assert.Equal(t, err.Error(), "Could not insert")
}

func TestCreateUserEmailError(t *testing.T) {
	pgdbs := NewPgDbService("mock", "")

	userEmail := ""
	userPwhash := "$2a$04$2a2qnoery/ULUw2WgKVd0OyeHhsHWINab9w9WTPoXqe8xY4PBrwXe"
	userName := "Test User"
	userApikey := GenerateApiKey()

	u, err := pgdbs.CreateUser(userEmail, userPwhash, userName, userApikey)

	if !assert.Nil(t, u, "User should be nil") {
		return
	}

	if !assert.NotNil(t, err, "Error should not be nil") {
		return
	}

	assert.Equal(t, err.Error(), "Email cannot be empty")
}

func TestCreateUser(t *testing.T) {
	pgdbs := NewPgDbService("mock", "")

	userNewId := 1
	userEmail := "test@example.com"
	userPwhash := "$2a$04$2a2qnoery/ULUw2WgKVd0OyeHhsHWINab9w9WTPoXqe8xY4PBrwXe"
	userName := "Test User"
	userApikey := GenerateApiKey()

	sqlmock.ExpectQuery(`INSERT INTO "user" \( email, pwhash, name, is_active, is_superuser, apikey \) VALUES \(\?, \?, \?, \?, \?, \?\) RETURNING id;`).
		WithArgs(userEmail, userPwhash, userName, true, false, userApikey).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userNewId))

	u, err := pgdbs.CreateUser(userEmail, userPwhash, userName, userApikey)

	if !assert.Nil(t, err, "Error should be nil") {
		return
	}

	if !assert.NotNil(t, u, "User should not be nil") {
		return
	}

	assert.Equal(t, u.Id, userNewId)
}
