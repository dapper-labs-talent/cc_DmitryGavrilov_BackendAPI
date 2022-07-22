package config

import (
	"os"

	"github.com/spf13/viper"
)

const (
	DefaultLogLevel      = 4     // logrus INFO log level
	DefaultJWTExpiration = 86400 // 1 day in seconds
	DefaultListenPort    = 8080
)

type Config struct {
	Server
	Database
	JWT
	Logging
}

type JWT struct {
	Expiration int    `mapstructure:"jwt_expiration"`
	Secret     string `mapstructure:"jwt_secret, required"`
}

type Database struct {
	Driver   string `mapstructure:"db_driver"`
	Host     string `mapstructure:"db_host"`
	DbName   string `mapstructure:"db_name"`
	Port     int    `mapstructure:"db_port"`
	User     string `mapstructure:"db_user"`
	Password string `mapstructure:"db_password"`
}

type Server struct {
	ListenPort int `mapstructure:"listen_port"`
	Timeout    int `mapstructure:"shutdown_timeout"`
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
