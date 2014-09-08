package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
			in:  `abc,def`,
			out: map[string]string{},
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

func TestParseCredentials(t *testing.T) {
	var parseCredentialsTests = []struct {
		in     string
		email  string
		apikey string
		err    error
	}{
		{
			in:     "test@example.com:abcdefg",
			email:  "test@example.com",
			apikey: "abcdefg",
			err:    nil,
		},
		{
			in:     `email="test@example.com", key="abcdefg"`,
			email:  "test@example.com",
			apikey: "abcdefg",
			err:    nil,
		},
		{
			in:     ``,
			email:  "",
			apikey: "",
			err:    SplitAuthError,
		},
		{
			in:     `a="A", b="B"`,
			email:  "",
			apikey: "",
			err:    AuthParseError,
		},
		{
			in:     "dGVzdEBleGFtcGxlLmNvbTphYmNkZWZn",
			email:  "test@example.com",
			apikey: "abcdefg",
			err:    nil,
		},
		{
			in:     `email="test@example.com"`,
			email:  "",
			apikey: "",
			err:    SplitAuthError,
		},
		{
			in:     `apikey="abcdefg"`,
			email:  "",
			apikey: "",
			err:    SplitAuthError,
		},
	}

	for i, test := range parseCredentialsTests {
		actualEmail, actualKey, actualErr := ParseCredentials(test.in)

		errMsg := fmt.Sprintf("Input %d: %q", i, test.in)
		assert.Equal(t, test.email, actualEmail, errMsg)
		assert.Equal(t, test.apikey, actualKey, errMsg)
		assert.Equal(t, test.err, actualErr, errMsg)
	}
}

func TestSplitAuthHeader(t *testing.T) {
	var splitAuthHeaderTests = []struct {
		in     string
		scheme string
		creds  string
		err    error
	}{
		{
			in:     "Apikey test@example.com:abcdefg",
			scheme: "Apikey",
			creds:  "test@example.com:abcdefg",
			err:    nil,
		},
		{
			in:     "",
			scheme: "",
			creds:  "",
			err:    AuthHeaderError,
		},
		{
			in:     "Apikey",
			scheme: "",
			creds:  "",
			err:    AuthHeaderError,
		},
		{
			in:     "    Apikey   ",
			scheme: "",
			creds:  "",
			err:    AuthHeaderError,
		},
		{
			in:     " Apikey ",
			scheme: "",
			creds:  "",
			err:    AuthHeaderError,
		},
		{
			in:     " Apikey  test@example.com abcdefg ",
			scheme: "Apikey",
			creds:  "test@example.com abcdefg",
			err:    nil,
		},
		{
			in:     `Apikey email="test@example.com", key="abcdefg"`,
			scheme: "Apikey",
			creds:  `email="test@example.com", key="abcdefg"`,
			err:    nil,
		},
		{
			in:     `api key email="test@example.com", key="abcdefg"`,
			scheme: "api",
			creds:  `key email="test@example.com", key="abcdefg"`,
			err:    nil,
		},
		{
			in:     "Apikey dGVzdEBleGFtcGxlLmNvbTphYmNkZWZn",
			scheme: "Apikey",
			creds:  "dGVzdEBleGFtcGxlLmNvbTphYmNkZWZn",
			err:    nil,
		},
	}

	for i, test := range splitAuthHeaderTests {
		actualScheme, actualCreds, actualErr := SplitAuthHeader(test.in)

		errMsg := fmt.Sprintf("Input %d: %q", i, test.in)
		assert.Equal(t, test.scheme, actualScheme, errMsg)
		assert.Equal(t, test.creds, actualCreds, errMsg)
		assert.Equal(t, test.err, actualErr, errMsg)
	}
}

func TestGetAuthHeader(t *testing.T) {
	var getAuthHeaderTests = []struct {
		in     http.Header
		scheme string
		creds  string
		err    error
	}{
		{
			in: http.Header{
				"Authorization": []string{"Apikey test@example.com:abcdefg"},
			},
			scheme: "Apikey",
			creds:  "test@example.com:abcdefg",
			err:    nil,
		},
		{
			in: http.Header{
				"Authentication": []string{"Apikey test@example.com:abcdefg"},
			},
			scheme: "",
			creds:  "",
			err:    AuthNotSetError,
		},
		{
			in: http.Header{
				"Authorization": []string{""},
			},
			scheme: "",
			creds:  "",
			err:    AuthHeaderError,
		},
		{
			in: http.Header{
				"Authorization": []string{"Basic test@example.com:asdf"},
			},
			scheme: "",
			creds:  "",
			err:    AuthTypeError,
		},
		{
			in: http.Header{
				"Authorization": []string{
					"Basic test@example.com:asdf",
					"Apikey test@example.com:abcdefg",
				},
			},
			scheme: "",
			creds:  "",
			err:    AuthTypeError,
		},
		{
			in: http.Header{
				"Authorization": []string{
					"Apikey test@example.com:abcdefg",
					"Basic test@example.com:asdf",
				},
			},
			scheme: "Apikey",
			creds:  "test@example.com:abcdefg",
			err:    nil,
		},
	}

	for i, test := range getAuthHeaderTests {
		actualScheme, actualCreds, actualErr := GetAuthHeader(test.in)

		errMsg := fmt.Sprintf("Input %d: %q", i, test.in)
		assert.Equal(t, test.scheme, actualScheme, errMsg)
		assert.Equal(t, test.creds, actualCreds, errMsg)
		assert.Equal(t, test.err, actualErr, errMsg)
	}
}

func TestUnauthorizedHeader(t *testing.T) {
	testRw := httptest.NewRecorder()
	UnauthorizedHeader(testRw)

	assert.Equal(t, testRw.Code, http.StatusUnauthorized)
	assert.Equal(t, testRw.Body.String(), ApiKeyRequiredError+"\n")

	actualHeader, ok := testRw.HeaderMap[http.CanonicalHeaderKey("WWW-Authenticate")]
	if !assert.True(t, ok, "WWW-Authenticate header should exist") {
		return
	}
	if !assert.Len(t, actualHeader, 1, "Header should have a value") {
		return
	}
	assert.Equal(t, actualHeader[0], "Apikey")
}
