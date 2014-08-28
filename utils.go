package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"github.com/lib/pq/hstore"
)

// Convert a given interface to JSON with indentation
func Jsonify(v interface{}) string {
	json_data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(json_data)
}

func GeneratePasswordHash(password string, cost int) string {
	pwhash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		panic(err)
	}
	return string(pwhash)
}

// Generate a unique string of length 40
// HMAC of a UUID4 with sha1
func GenerateApiKey() string {
	mac := hmac.New(sha1.New, uuid.NewRandom())
	return hex.EncodeToString(mac.Sum(nil))
}

func MapToHstore(m map[string]string, h *hstore.Hstore) {
	h.Map = make(map[string]sql.NullString)

	for key, val := range m {
		h.Map[key] = sql.NullString{val, true}
	}
}

func HstoreToMap(h *hstore.Hstore) map[string]string {
	m := make(map[string]string)
	if h.Map != nil {
		for key, val := range h.Map {
			if val.Valid {
				m[key] = val.String
			}
		}
	}

	return m
}
