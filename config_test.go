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
				DbConf: dbConfig{
					User: "test1",
				},
			},
		},
		{
			in: `{db: {user: test1, password: test2}}`,
			out: appConfig{
				DbConf: dbConfig{
					User:     "test1",
					Password: "test2",
				},
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
				DbConf: dbConfig{
					User:     "test1",
					Password: "test2",
				},
				ListenConf: listenConfig{
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
				DbConf: dbConfig{
					Type:     "mysql",
					Host:     "dbhost",
					Port:     7654,
					User:     "test1",
					Password: "test2",
					DbName:   "tester3",
					SslMode:  "verify-full",
				},
				ListenConf: listenConfig{
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
				DbConf: dbConfig{
					Type: "postgres",
				},
			},
			out: "postgres",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					Type: "mysql",
				},
			},
			out: "mysql",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					Type: "mock",
				},
			},
			out: "mock",
		},
	}

	for _, test := range validateTests {
		assert.Equal(t, test.out, test.in.DbType())
	}
}

func TestReadConfigError(t *testing.T) {
	invalidString := "{{{"

	actualOut, err := ReadConfig([]byte(invalidString))
	assert.Nil(t, actualOut)
	assert.NotNil(t, err)
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
				DbConf: dbConfig{
					Host: "test1",
				},
			},
			out: "host=test1 port=5432 application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					Port: 6543,
				},
			},
			out: "host=localhost port=6543 application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					Host: "test1",
					Port: 7654,
				},
			},
			out: "host=test1 port=7654 application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					User: "test1",
				},
			},
			out: "host=localhost port=5432 user=test1 application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					User:     "test1",
					Password: "test2",
				},
			},
			out: "host=localhost port=5432 user=test1 password=test2 application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					Password: "test2",
				},
			},
			out: "host=localhost port=5432 password=test2 application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					DbName: "testdb",
				},
			},
			out: "host=localhost port=5432 name=testdb application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					User:   "test1",
					DbName: "testdb",
				},
			},
			out: "host=localhost port=5432 user=test1 name=testdb application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					User:     "test1",
					Password: "test2",
					DbName:   "testdb",
				},
			},
			out: "host=localhost port=5432 user=test1 password=test2 name=testdb application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					Password: "test2",
					DbName:   "testdb",
				},
			},
			out: "host=localhost port=5432 password=test2 name=testdb application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					SslMode: "verify-full",
				},
			},
			out: "host=localhost port=5432 sslmode=verify-full application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					DbName:  "testdb",
					SslMode: "verify-full",
				},
			},
			out: "host=localhost port=5432 name=testdb sslmode=verify-full application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					User:     "test1",
					Password: "test2",
					DbName:   "testdb",
					SslMode:  "verify-full",
				},
			},
			out: "host=localhost port=5432 user=test1 password=test2 name=testdb sslmode=verify-full application_name=people-go",
		},
	}

	for _, test := range validateTests {
		assert.Equal(t, test.out, test.in.DbCreds())
	}
}
