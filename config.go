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

type dbConfig struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
	DbName   string `yaml="name"`
}

type listenConfig struct {
	Host string
	Port int
}

type appConfig struct {
	DbConf     dbConfig     `yaml="db"`
	ListenConf listenConfig `yaml="listen"`
}

func (ac *appConfig) DbType() string {
	return ""
}

func (ac *appConfig) DbCreds() string {
	return ""
}

func (ac *appConfig) ListenAddr() string {
	return ""
}

