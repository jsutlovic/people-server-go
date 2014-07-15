package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitAuth(t *testing.T) {
	var splitAuthTests = []struct {
		in    string
		email string
		key   string
		err   error
	}{
		{
			in:    "test@example.com:abcdefg",
			email: "test@example.com",
			key:   "abcdefg",
			err:   nil,
		},
		{
			in:    "dGVzdEBleGFtcGxlLmNvbTphYmNkZWZn",
			email: "test@example.com",
			key:   "abcdefg",
			err:   nil,
		},
		{
			in:    "test@example.com",
			email: "",
			key:   "",
			err:   SplitAuthError,
		},
		{
			in:    "abcdefg",
			email: "",
			key:   "",
			err:   SplitAuthError,
		},
		{
			in:    "dGVzdEBleGFtcGxlLmNvbQ==",
			email: "",
			key:   "",
			err:   SplitAuthError,
		},
		{
			in:    "YWJjZGVmZw==",
			email: "",
			key:   "",
			err:   SplitAuthError,
		},
		{
			in:    "test@example.com:abcdefg:hijk",
			email: "test@example.com",
			key:   "abcdefg:hijk",
			err:   nil,
		},
		{
			in:    "dGVzdEBleGFtcGxlLmNvbTphYmNkZWZnOmhpams=",
			email: "test@example.com",
			key:   "abcdefg:hijk",
			err:   nil,
		},
		{
			in:    "",
			email: "",
			key:   "",
			err:   SplitAuthError,
		},
	}

	for i, test := range splitAuthTests {
		actualEmail, actualKey, actualErr := SplitAuth(test.in)
		assert.Equal(t, test.email, actualEmail, "Input %d: %q", i, test.in)
		assert.Equal(t, test.key, actualKey, "Input %d: %q", i, test.in)
		assert.Equal(t, test.err, actualErr, "Input %d: %q", i, test.in)
	}
}

func TestSplitFields(t *testing.T) {
	var splitFieldsTests = []struct {
		in  string
		out map[string]string
	}{
		{
			in: `email="test@example.com", key="abcdefg"`,
			out: map[string]string{
				"email": "test@example.com",
				"key":   "abcdefg",
			},
		},
		{
			in: `email="test@example.com", key="abcdefg",`,
			out: map[string]string{
				"email": "test@example.com",
				"key":   "abcdefg",
			},
		},
		{
			in:  ``,
			out: map[string]string{},
		},
		{
			in: `a=abc,b=def`,
			out: map[string]string{
				"a": "abc",
				"b": "def",
			},
		},
		{
			in: `a = "abc"  ,  b = "def " `,
			out: map[string]string{
				"a": "abc",
				"b": "def ",
			},
		},
		{
			in: `a = " a b c ", b = "  d  e  f  ", c="ghi"`,
			out: map[string]string{
				"a": " a b c ",
				"b": "  d  e  f  ",
				"c": "ghi",
			},
		},
		{
			in:  `test@example.com:abcdefg`,
			out: map[string]string{},
		},
		{
			in:  `dGVzdEBleGFtcGxlLmNvbTphYmNkZWZn`,
			out: map[string]string{},
		},
		{
			in: `a="abc" b="def"`,
			out: map[string]string{
				"a": `abc" b="def`,
			},
		},
		{
			in: `a="abc"`,
			out: map[string]string{
				"a": "abc",
			},
		},
		{
			in: `"a"="abc"`,
			out: map[string]string{
				`"a"`: "abc",
			},
		},
		{
			in: ` a =   "abc"  `,
			out: map[string]string{
				"a": "abc",
			},
		},
		{
			in: `a=abc, b="def"`,
			out: map[string]string{
				"a": "abc",
				"b": "def",
			},
		},
		{
			in: `a="abc, b=def"`,
			out: map[string]string{
				"a": "abc",
				"b": "def",
			},
		},
		{
			in: `a=abc", "b = def`,
			out: map[string]string{
				`a`:  "abc",
				`"b`: "def",
			},
		},
	}

	for i, test := range splitFieldsTests {
		actual := SplitFields(test.in)
		assert.Equal(t, test.out, actual, "Input %d: %q", i, test.in)
	}
}
