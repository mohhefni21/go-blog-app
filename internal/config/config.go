package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	AppConfig AppConfig
	DBconfig  DBconfig
}

type AppConfig struct {
	AppName string `env:"APPNAME,required"`
	AppPort string `env:"APPPORT,required"`
}

type DBconfig struct {
	DBHost       string `env:"DBHOST,required"`
	DBPort       string `env:"DBPORT,required"`
	DBusername   string `env:"DBUSERNAME,required"`
	DBPassword   string `env:"DBPASSWORD,required"`
	DBName       string `env:"DBNAME,required"`
	DBPoolConfig DBPoolConfig
}

type DBPoolConfig struct {
	MaxIdleConnection      uint8 `env:"MAXIDLECONNECTION,required"`
	MaxOpenConnection      uint8 `env:"MAXOPENCONNECTION,required"`
	MaxLifetimeConnection  uint8 `env:"MAXLIFETIMECONNECTION,required"`
	MaxIdleTimeConnetction uint8 `env:"MAXIDLETIMECONNECTION,required"`
}

var Cfg Config

func LoadConfig(path string) (err error) {
	err = godotenv.Load(path)
	if err != nil {
		return
	}

	err = env.Parse(&Cfg)

	return
}
