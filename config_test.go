package main

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

const (
	testConfigFile            = "config.yml.example"
	testNonexistentConfigFile = "nonexistent.yml"
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
  address: 0.0.0.0
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
					Address: "0.0.0.0",
					Port:    4321,
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

func TestReadConfigError(t *testing.T) {
	invalidString := "{{{"

	actualOut, err := ReadConfig([]byte(invalidString))
	assert.Nil(t, actualOut)
	assert.NotNil(t, err)
}

func TestReadConfigFileParses(t *testing.T) {
	expected := appConfig{
		DbConf: dbConfig{
			Type:     "postgres",
			Host:     "dbhost",
			Port:     5432,
			User:     "people-user",
			Password: "people-pw",
			DbName:   "people-db",
			SslMode:  "disable",
		},
		ListenConf: listenConfig{
			Address: "0.0.0.0",
			Port:    3001,
		},
	}
	actualOut, err := ReadConfigFile(testConfigFile)
	assert.Nil(t, err)
	assert.Equal(t, &expected, actualOut)
}

func TestReadConfigFileError(t *testing.T) {
	actualOut, err := ReadConfigFile(testNonexistentConfigFile)
	assert.Nil(t, actualOut)
	assert.NotNil(t, err)
}

func TestMustReadConfigFileParses(t *testing.T) {
	expected := appConfig{
		DbConf: dbConfig{
			Type:     "postgres",
			Host:     "dbhost",
			Port:     5432,
			User:     "people-user",
			Password: "people-pw",
			DbName:   "people-db",
			SslMode:  "disable",
		},
		ListenConf: listenConfig{
			Address: "0.0.0.0",
			Port:    3001,
		},
	}
	actualOut := MustReadConfigFile(testConfigFile)
	assert.Equal(t, &expected, actualOut)
}

func TestMustReadConfigFilePanics(t *testing.T) {
	assert.Panics(t, func() {
		_ = MustReadConfigFile(testNonexistentConfigFile)
	})
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
			out: "host=localhost port=5432 dbname=testdb application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					User:   "test1",
					DbName: "testdb",
				},
			},
			out: "host=localhost port=5432 user=test1 dbname=testdb application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					User:     "test1",
					Password: "test2",
					DbName:   "testdb",
				},
			},
			out: "host=localhost port=5432 user=test1 password=test2 dbname=testdb application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					Password: "test2",
					DbName:   "testdb",
				},
			},
			out: "host=localhost port=5432 password=test2 dbname=testdb application_name=people-go",
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
			out: "host=localhost port=5432 dbname=testdb sslmode=verify-full application_name=people-go",
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
			out: "host=localhost port=5432 user=test1 password=test2 dbname=testdb sslmode=verify-full application_name=people-go",
		},
		{
			in: appConfig{
				DbConf: dbConfig{
					Type:     "mysql",
					Host:     "testhost",
					Port:     3306,
					User:     "test1",
					Password: "test2",
					DbName:   "testdb",
					SslMode:  "verify-full",
				},
			},
			out: "host=testhost port=3306 user=test1 password=test2 dbname=testdb sslmode=verify-full application_name=people-go",
		},
	}

	for _, test := range validateTests {
		assert.Equal(t, test.out, test.in.DbCreds())
	}
}

func TestAppConfigListenerAddr(t *testing.T) {
	validateTests := []struct {
		in  appConfig
		out string
	}{
		{
			in:  appConfig{},
			out: "127.0.0.1:3000",
		},
		{
			in: appConfig{
				ListenConf: listenConfig{
					Address: "0.0.0.0",
				},
			},
			out: "0.0.0.0:3000",
		},
		{
			in: appConfig{
				ListenConf: listenConfig{
					Port: 4040,
				},
			},
			out: "127.0.0.1:4040",
		},
		{
			in: appConfig{
				ListenConf: listenConfig{
					Address: "0.0.0.0",
					Port:    4040,
				},
			},
			out: "0.0.0.0:4040",
		},
		{
			in: appConfig{
				ListenConf: listenConfig{
					Address: "::1",
					Ipv6:    true,
				},
			},
			out: "[::1]:3000",
		},
		{
			in: appConfig{
				ListenConf: listenConfig{
					Address: "::1",
					Port:    4040,
					Ipv6:    true,
				},
			},
			out: "[::1]:4040",
		},
		{
			in: appConfig{
				ListenConf: listenConfig{
					Address: "/tmp/people-test.sock",
				},
			},
			out: "/tmp/people-test.sock",
		},
		{
			in: appConfig{
				ListenConf: listenConfig{
					Address: "/tmp/people-test.sock",
					Port:    4040,
				},
			},
			out: "/tmp/people-test.sock",
		},
		{
			in: appConfig{
				ListenConf: listenConfig{
					Address: "/tmp/people-test.sock",
					Ipv6:    true,
				},
			},
			out: "/tmp/people-test.sock",
		},
	}

	for i, test := range validateTests {
		msg := strconv.Itoa(i)
		safe := assert.NotPanics(t, func() {
			l := test.in.Listener()
			l.Close()
		}, msg)

		if safe {
			l := test.in.Listener()
			if assert.NotNil(t, l, msg) {
				assert.Equal(t, test.out, l.Addr().String(), msg)
				l.Close()
			}
		}
	}
}

func TestAppConfigListenerPanics(t *testing.T) {
	invalidTests := []appConfig{
		appConfig{
			ListenConf: listenConfig{
				Address: "127.0.0.1",
				Ipv6:    true,
			},
		},
		appConfig{
			ListenConf: listenConfig{
				Address: "::1",
			},
		},
	}

	for i, test := range invalidTests {
		assert.Panics(t, func() {
			test.Listener()
		}, strconv.Itoa(i))
	}
}
