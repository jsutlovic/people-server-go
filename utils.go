package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
)

// Convert a given interface to JSON with indentation
func Jsonify(v interface{}) string {
	json_data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(json_data)
}

func GeneratePasswordHash(password string, cost int) (string, error) {
	pwhash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(pwhash), nil
}

// Generate a unique string of length 40
// HMAC of a UUID4 with sha1
func GenerateApiKey() string {
	mac := hmac.New(sha1.New, uuid.NewRandom())
	return hex.EncodeToString(mac.Sum(nil))
}
