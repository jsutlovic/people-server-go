package main

// Defaults
const (
	DefaultHost   = "127.0.0.1"
	DefaultPort   = 3000
	DefaultDbType = "postgres"
	DefaultDbHost = "localhost"
	DefaultDbPort = 5432
)

type Config interface {
	DbType() string
	DbCreds() string
	ListenAddr() string
}
