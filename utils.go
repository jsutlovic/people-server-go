package main

import (
	"code.google.com/p/go.crypto/bcrypt"
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

func CreatePassword(password string, cost int) (string, error) {
	pwhash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(pwhash), nil
}
