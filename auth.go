package main

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

const (
	AuthHeaderKey string = "Authorization"
	EmailField    string = "email"
	KeyField      string = "key"
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
Parse an http.Request
*/
func ParseRequestHeaders(req *http.Request) (authParams *AuthParams, err error) {
	h, ok := req.Header[http.CanonicalHeaderKey(AuthHeaderKey)]
	if !ok || len(h) == 0 {
		return nil, errors.New("Authorization header not set")
	}
	scheme, creds, err := SplitHeader(h[0])
	if err != nil {
		return nil, err
	}
	if strings.ToLower(scheme) != "apikey" {
		return nil, errors.New("Authorization type is not Apikey")
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

func SplitHeader(h string) (scheme, credentials string, err error) {
	parts := strings.SplitN(h, " ", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
	}
	return "", "", errors.New("Could not parse Authorization header")
}

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
		return "", "", errors.New("Could not parse authentication")
	}
	return email, apikey, nil
}

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
		return "", "", errors.New("Could not parse auth string")
	}

	email = strings.TrimSpace(fields[0])
	key = strings.TrimSpace(fields[1])
	return email, key, nil
}
