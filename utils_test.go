package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePasswordCostTooLarge(t *testing.T) {
	cost := 32
	_, err := CreatePassword("asdf", cost)

	if !assert.NotNil(t, err, "Cost should error") {
		t.Logf("Cost: %d", cost)
	}
}

func TestCreatePassword(t *testing.T) {
	passwords := []string{
		"",
		"a",
		"ab",
		"abc",
		"abcd",
		"abcde",
		"abcdef",
		"abcdefg",
		"abcdefgh",
		"abcdefghi",
		"abcdefghij",
		"abcdefghijk",
		"abcdefghijkl",
		"abcdefghijklm",
		"abcdefghijklmn",
	}

	user := new(User)

	for _, password := range passwords {
		hashed, err := CreatePassword(password, 4)

		if !assert.Nil(t, err, "Password should not error") {
			break
		}

		user.Pwhash = hashed

		if !assert.True(t, user.CheckPassword(password)) {
			break
		}
	}
}

func TestCreateApiKey(t *testing.T) {
	// Create a number of keys
	loops := 10
	keys := make([]string, loops)

	for i := 0; i < loops; i++ {
		key := GenerateApiKey()

		// Key length should be 40
		assert.Len(t, key, 40)

		// Previous keys should not match
		for j := 0; j < i; j++ {
			if !assert.NotEqual(t, keys[j], key) {
				break
			}
		}
		keys[i] = key
	}
}
