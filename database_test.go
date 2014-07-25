package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func newTestUser() User {
	user := User{
		1,
		"test@example.com",
		"",
		"Test User",
		true,
		false,
		"",
	}

	return user
}

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
