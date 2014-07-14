package main

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

// Request fields
const (
	AuthHeaderKey string = "Authorization"
	EmailField    string = "email"
	KeyField      string = "key"
)

// Errors
var (
	SplitAuthError  error = errors.New("Could not parse auth string")
	AuthParseError  error = errors.New("Could not parse authentication")
	AuthHeaderError error = errors.New("Could not parse Authorization header")
	AuthNotSetError error = errors.New("Authorization header not set")
	AuthTypeError   error = errors.New("Authorization type is not Apikey")
)

type AuthParams struct {
	Authed bool
	Email  string
	Apikey string
}

func Authorize(req *http.Request) (authParams *AuthParams, err error) {
	authParams, err = ParseRequestHeaders(req)

	if err != nil {
		return nil, err
	}

	return authParams, nil
}

/*
Get authorization parameters from an http.Request

Looks for the HTTP Authorization header and parses it
*/
func ParseRequestHeaders(req *http.Request) (authParams *AuthParams, err error) {
	h, ok := req.Header[http.CanonicalHeaderKey(AuthHeaderKey)]
	if !ok || len(h) == 0 {
		return nil, AuthNotSetError
	}
	scheme, creds, err := SplitAuthHeader(h[0])
	if err != nil {
		return nil, err
	}
	if strings.ToLower(scheme) != "apikey" {
		return nil, AuthTypeError
	}
	email, apikey, err := ParseCredentials(creds)
	if err != nil {
		return nil, err
	}

	authParams = new(AuthParams)
	authParams.Authed = true
	authParams.Email = email
	authParams.Apikey = apikey

	return authParams, nil
}

// Split the header into auth scheme and credentials
func SplitAuthHeader(h string) (scheme, credentials string, err error) {
	parts := strings.SplitN(h, " ", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
	}
	return "", "", AuthHeaderError
}

// Parse auth credentials into email and apikey fields
func ParseCredentials(creds string) (email, apikey string, err error) {
	fields := SplitFields(creds)

	// We didn't find the proper keys
	if len(fields) < 2 {
		// Try colon-separated
		return SplitAuth(creds)
	}

	email, emailOk := fields[EmailField]
	apikey, apikeyOk := fields[KeyField]

	if !(emailOk && apikeyOk) {
		return "", "", AuthParseError
	}
	return email, apikey, nil
}

/*
Split key/value HTTP Authorization parameters

Parses: "email=\"test@example.com\", key=\"abcdefg\""
into {email: "test@example.com" key: "abcdefg"}
*/
func SplitFields(text string) (fields map[string]string) {
	fields = make(map[string]string)

	raw_fields := strings.Split(text, ",")

	for _, raw_field := range raw_fields {
		splitField := strings.SplitN(raw_field, "=", 2)
		if len(splitField) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(splitField[0]))
		val := strings.Trim(splitField[1], "\" ")
		fields[key] = val
	}

	return fields
}

/*
Split Authorization credentials into email and apikey

Parses "test@example.com:abcdefg"
into "test@example.com", "abcdefg"

Supports base64 encoded strings
*/
func SplitAuth(raw_text string) (email, key string, err error) {
	var text string
	decoded, err := base64.StdEncoding.DecodeString(raw_text)

	if err != nil {
		text = raw_text
	} else {
		text = string(decoded)
	}

	fields := strings.SplitN(text, ":", 2)
	if len(fields) != 2 {
		return "", "", SplitAuthError
	}

	email = strings.TrimSpace(fields[0])
	key = strings.TrimSpace(fields[1])
	return email, key, nil
}
