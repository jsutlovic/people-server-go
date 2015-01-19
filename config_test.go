package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfigParses(t *testing.T) {
	validateTests := []struct {
		in  string
		out appConfig
	}{
		{
			in:  ``,
			out: appConfig{},
		},
		{
			in: `{db: {user: test1}}`,
			out: appConfig{
				dbConfig{
					User: "test1",
				},
				listenConfig{},
			},
		},
		{
			in: `{db: {user: test1, password: test2}}`,
			out: appConfig{
				dbConfig{
					User:     "test1",
					Password: "test2",
				},
				listenConfig{},
			},
		},
		{
			in: `
---
db:
  user: test1
  password: test2

listen:
  port: 4321
`,
			out: appConfig{
				dbConfig{
					User:     "test1",
					Password: "test2",
				},
				listenConfig{
					Port: 4321,
				},
			},
		},
		{
			in: `
---
db:
  type: mysql
  host: dbhost
  port: 7654
  user: test1
  password: test2
  name: tester3
  sslmode: verify-full

listen:
  host: 0.0.0.0
  port: 4321
`,
			out: appConfig{
				dbConfig{
					Type:     "mysql",
					Host:     "dbhost",
					Port:     7654,
					User:     "test1",
					Password: "test2",
					DbName:   "tester3",
					SslMode:  "verify-full",
				},
				listenConfig{
					Host: "0.0.0.0",
					Port: 4321,
				},
			},
		},
	}

	for _, test := range validateTests {
		actualOut, err := ReadConfig([]byte(test.in))
		assert.Nil(t, err)
		assert.Equal(t, &test.out, actualOut)
	}
}

func TestAppConfigDbCreds(t *testing.T) {
	validateTests := []struct {
		in  appConfig
		out string
	}{
		{
			in:  appConfig{},
			out: "host=localhost port=5432 application_name=people-go",
		},
		{
			in: appConfig{
				dbConfig{
					User: "test1",
				},
				listenConfig{},
			},
			out: "host=localhost port=5432 user=test1 application_name=people-go",
		},
	}

	for _, test := range validateTests {
		assert.Equal(t, test.out, test.in.DbCreds())
	}
}
