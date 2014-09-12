package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq/hstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonifyError(t *testing.T) {
	assert.Panics(t, func() {
		Jsonify(func() {})
	})

	assert.Panics(t, func() {
		var test complex128 = 0
		Jsonify(test)
	})
}

func TestCreatePasswordCostTooLarge(t *testing.T) {
	cost := 32
	assert.Panics(t, func() {
		GeneratePasswordHash("asdf", cost)
	})
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
		hashed := GeneratePasswordHash(password, 4)

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

func TestMapToHstore(t *testing.T) {
	var mapToHstoreTests = []struct {
		in  map[string]string
		out hstore.Hstore
	}{
		{
			in: map[string]string{
				"a": "b",
				"c": "d",
			},
			out: hstore.Hstore{map[string]sql.NullString{
				"a": sql.NullString{"b", true},
				"c": sql.NullString{"d", true},
			}},
		},
		{
			in: map[string]string{
				"a": "",
				"b": "",
			},
			out: hstore.Hstore{map[string]sql.NullString{
				"a": sql.NullString{"", true},
				"b": sql.NullString{"", true},
			}},
		},
		{
			in:  map[string]string{},
			out: hstore.Hstore{map[string]sql.NullString{}},
		},
		{
			in:  nil,
			out: hstore.Hstore{map[string]sql.NullString{}},
		},
	}

	for _, test := range mapToHstoreTests {
		h := hstore.Hstore{}
		MapToHstore(test.in, &h)
		assert.Equal(t, h, test.out)
	}
}

func TestHstoreToMap(t *testing.T) {
	var hstoreToMapTests = []struct {
		in  hstore.Hstore
		out map[string]string
	}{
		{
			in: hstore.Hstore{map[string]sql.NullString{
				"a": sql.NullString{"b", true},
				"c": sql.NullString{"d", true},
			}},
			out: map[string]string{
				"a": "b",
				"c": "d",
			},
		},
		{
			in: hstore.Hstore{map[string]sql.NullString{
				"a": sql.NullString{"", true},
				"b": sql.NullString{"", true},
			}},
			out: map[string]string{
				"a": "",
				"b": "",
			},
		},
		{
			in: hstore.Hstore{map[string]sql.NullString{
				"a": sql.NullString{"", false},
				"b": sql.NullString{"asdfasdf", false},
			}},
			out: map[string]string{},
		},
		{
			in:  hstore.Hstore{nil},
			out: map[string]string{},
		},
	}

	for _, test := range hstoreToMapTests {
		assert.Equal(t, HstoreToMap(&test.in), test.out)
	}
}

func TestValidateEmail(t *testing.T) {
	emailTests := []struct {
		in  string
		out bool
	}{
		// Invalid emails
		{
			in:  "",
			out: false,
		},
		{
			in:  "test",
			out: false,
		},
		{
			in:  "test@",
			out: false,
		},
		{
			in:  "test@example",
			out: false,
		},
		{
			in:  "test@example.",
			out: false,
		},
		{
			in:  "test.example@example",
			out: false,
		},
		{
			in:  "test example@example.com ",
			out: false,
		},
		{
			in:  "test.example@ example.com",
			out: false,
		},
		{
			in:  "test.example@ example.com ",
			out: false,
		},
		{
			in:  " ",
			out: false,
		},
		{
			in:  " test @example.com",
			out: false,
		},
		{
			in:  " test @ example com",
			out: false,
		},
		{
			in:  "test@test@example.com",
			out: false,
		},
		{
			in:  "test.test@example@example.co.uk",
			out: false,
		},

		// Valid emails
		{
			in:  "test@example.com",
			out: true,
		},
		{
			in:  "test.example@example.com",
			out: true,
		},
		{
			in:  "test.example@example.co.uk",
			out: true,
		},
		{
			in:  "test@example.co.uk",
			out: true,
		},
	}

	for i, test := range emailTests {
		assert.Equal(
			t,
			ValidateEmail(test.in),
			test.out,
			fmt.Sprintf("Test %d: %#v", i+1, test.in))
	}
}
