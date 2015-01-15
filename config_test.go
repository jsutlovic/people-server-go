package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigParses(t *testing.T) {
	validateTests := []struct {
		in  string
		out Config
	}{
		{
			// Defaults should be set if input is empty
			in: ``,
			out: &mockConfig{
				"postgres",
				"host=localhost port=5432 application_name=people-go",
				"127.0.0.1:3000",
			},
		},
	}

	for _, test := range validateTests {
		actualOut, err := ReadConfig([]byte(test.in))
		assert.Nil(t, err)
		assert.Equal(t, actualOut.DbType(), test.out.DbType())
		assert.Equal(t, actualOut.DbCreds(), test.out.DbCreds())
		assert.Equal(t, actualOut.ListenAddr(), test.out.ListenAddr())
	}
}
