package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	AppConfig   AppConfig
	DBconfig    DBconfig
	AuthConfig  AuthConfig
	OAuthConfig OAuthConfig
}

type AppConfig struct {
	AppName string `env:"APPNAME,required"`
	AppPort string `env:"APPPORT,required"`
	BaseUrl string `env:"BASE_URL,required"`
}

type DBconfig struct {
	DBHost       string `env:"DBHOST,required"`
	DBPort       string `env:"DBPORT,required"`
	DBusername   string `env:"DBUSERNAME,required"`
	DBPassword   string `env:"DBPASSWORD,required"`
	DBName       string `env:"DBNAME,required"`
	DBPoolConfig DBPoolConfig
}

type OAuthConfig struct {
	GoogleClientId     string `env:"GOOGLECLIENTID,required"`
	GoogleClientSecret string `env:"GOOGLECLIENTSECRET,required"`
	ClientCallbackUrl  string `env:"CLIENTCALLBACKURL,required"`
	GoogleStateToken   string `env:"GOOGLESTATETOKEN,required"`
}

type DBPoolConfig struct {
	MaxIdleConnection      uint8 `env:"MAXIDLECONNECTION,required"`
	MaxOpenConnection      uint8 `env:"MAXOPENCONNECTION,required"`
	MaxLifetimeConnection  uint8 `env:"MAXLIFETIMECONNECTION,required"`
	MaxIdleTimeConnetction uint8 `env:"MAXIDLETIMECONNECTION,required"`
}

type AuthConfig struct {
	AccessTokenKey         string `env:"ACCESSTOKENKEY,required"`
	RefreshTokenKey        string `env:"REFRESHTOKENKEY,required"`
	AccessTokenExpiration  int    `env:"ACCESSTOKENEXPIRATION,required"`
	RefreshTokenExpiration int    `env:"REFRESHTOKENEXPIRATION,required"`
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
