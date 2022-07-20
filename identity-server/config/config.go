package config

import (
	"os"

	"github.com/spf13/viper"
)

const (
	DefaultLogLevel      = 3     // INFO log level
	DefaultJWTExpiration = 86400 // 1 day in seconds
	DefaultListenPort    = 8080
)

type Config struct {
	Server
	Postgres
	JWT
	Logging
}

type JWT struct {
	Exp        int    `mapstructure:"jwt_expiration"`
	SignSecret string `mapstructure:"jwt_secret, required"`
}

type Postgres struct {
	Host     string `mapstructure:"postgres_host"`
	Database string `mapstructure:"postgres_database"`
	Port     int    `mapstructure:"postgres_port"`
	User     string `mapstructure:"postgres_user"`
	Password string `mapstructure:"postgres_password"`
}

type Server struct {
	ListenPort int `mapstructure:"listen_port"`
}

type Logging struct {
	LogLevel int `mapstructure:"log_level"`
}

func LoadConfigWithPath(path string) (*Config, error) {
	if _, err := os.Stat(path); err == nil {
		viper.SetConfigFile(path)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("./etc/")
	viper.AddConfigPath("/usr/dapper-labs/etc/")

	viper.SetDefault("log_level", DefaultLogLevel)
	viper.SetDefault("jwt_expiration", DefaultJWTExpiration)
	viper.SetDefault("listen_port", DefaultListenPort)
}
