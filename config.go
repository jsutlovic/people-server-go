package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

// Defaults
const (
	DefaultAddress    = "127.0.0.1"
	DefaultPort       = 3000
	DefaultDbType     = "postgres"
	DefaultDbHost     = "localhost"
	DefaultDbPort     = 5432
	DefaultSslMode    = "disable"
	DbApplicationName = "people-go"
)

const (
	KeyValTemplate     = "%s=%v"
	ListenAddrTemplate = "%s:%d"
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
	DbName   string `yaml:"name"`
	SslMode  string `yaml:"sslmode"`
}

type listenConfig struct {
	Address string
	Port    int
}

type appConfig struct {
	DbConf     dbConfig     `yaml:"db"`
	ListenConf listenConfig `yaml:"listen"`
}

func (ac *appConfig) DbType() string {
	return defaultString(ac.DbConf.Type, DefaultDbType)
}

func (ac *appConfig) DbCreds() string {
	configStrings := []string{}

	hostStr := fmt.Sprintf(
		KeyValTemplate, "host", defaultString(ac.DbConf.Host, DefaultDbHost))

	portStr := fmt.Sprintf(
		KeyValTemplate, "port", defaultInt(ac.DbConf.Port, DefaultDbPort))

	configStrings = append(configStrings, hostStr, portStr)

	if ac.DbConf.User != "" {
		configStrings = append(configStrings,
			fmt.Sprintf(KeyValTemplate, "user", ac.DbConf.User))
	}

	if ac.DbConf.Password != "" {
		configStrings = append(configStrings,
			fmt.Sprintf(KeyValTemplate, "password", ac.DbConf.Password))
	}

	if ac.DbConf.DbName != "" {
		configStrings = append(configStrings,
			fmt.Sprintf(KeyValTemplate, "dbname", ac.DbConf.DbName))
	}

	if ac.DbConf.SslMode != "" {
		configStrings = append(configStrings,
			fmt.Sprintf(KeyValTemplate, "sslmode", ac.DbConf.SslMode))
	}

	appNameStr := fmt.Sprintf(
		KeyValTemplate, "application_name", DbApplicationName)

	configStrings = append(configStrings, appNameStr)

	return strings.Join(configStrings, " ")
}

func (ac *appConfig) ListenAddr() string {
	return fmt.Sprintf(
		ListenAddrTemplate,
		defaultString(ac.ListenConf.Address, DefaultAddress),
		defaultInt(ac.ListenConf.Port, DefaultPort))
}

// Always returns a string. If chk is empty, returns def
func defaultString(chk, def string) string {
	if chk == "" {
		return def
	}
	return chk
}

func defaultInt(chk, def int) int {
	if chk == 0 {
		return def
	}
	return chk
}

func MustReadConfigFile(filename string) *appConfig {
	config, err := ReadConfigFile(filename)
	if err != nil {
		panic(err)
	}
	return config
}

func ReadConfigFile(filename string) (*appConfig, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ReadConfig(b)
}

func ReadConfig(b []byte) (*appConfig, error) {
	config := &appConfig{}

	err := yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
