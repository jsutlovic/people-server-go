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

func TestAppConfigDbTyoe(t *testing.T) {
	validateTests := []struct {
		in  appConfig
		out string
	}{
		{
			in:  appConfig{},
			out: "postgres",
		},
		{
			in: appConfig{
				dbConfig{
					Type: "postgres",
				},
				listenConfig{},
			},
			out: "postgres",
		},
		{
			in: appConfig{
				dbConfig{
					Type: "mysql",
				},
				listenConfig{},
			},
			out: "mysql",
		},
		{
			in: appConfig{
				dbConfig{
					Type: "mock",
				},
				listenConfig{},
			},
			out: "mock",
		},
	}

	for _, test := range validateTests {
		assert.Equal(t, test.out, test.in.DbType())
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
					Host: "test1",
				},
				listenConfig{},
			},
			out: "host=test1 port=5432 application_name=people-go",
		},
		{
			in: appConfig{
				dbConfig{
					Port: 6543,
				},
				listenConfig{},
			},
			out: "host=localhost port=6543 application_name=people-go",
		},
		{
			in: appConfig{
				dbConfig{
					Host: "test1",
					Port: 7654,
				},
				listenConfig{},
			},
			out: "host=test1 port=7654 application_name=people-go",
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
		{
			in: appConfig{
				dbConfig{
					User:     "test1",
					Password: "test2",
				},
				listenConfig{},
			},
			out: "host=localhost port=5432 user=test1 password=test2 application_name=people-go",
		},
		{
			in: appConfig{
				dbConfig{
					Password: "test2",
				},
				listenConfig{},
			},
			out: "host=localhost port=5432 password=test2 application_name=people-go",
		},
		{
			in: appConfig{
				dbConfig{
					DbName: "testdb",
				},
				listenConfig{},
			},
			out: "host=localhost port=5432 name=testdb application_name=people-go",
		},
		{
			in: appConfig{
				dbConfig{
					User:   "test1",
					DbName: "testdb",
				},
				listenConfig{},
			},
			out: "host=localhost port=5432 user=test1 name=testdb application_name=people-go",
		},
		{
			in: appConfig{
				dbConfig{
					User:     "test1",
					Password: "test2",
					DbName:   "testdb",
				},
				listenConfig{},
			},
			out: "host=localhost port=5432 user=test1 password=test2 name=testdb application_name=people-go",
		},
	}

	for _, test := range validateTests {
		assert.Equal(t, test.out, test.in.DbCreds())
	}
}
