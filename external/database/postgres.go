package database

import (
	"fmt"
	"mohhefni/go-blog-app/internal/config"
	"time"

	"github.com/jmoiron/sqlx"
	// PostgreSQL driver
	_ "github.com/lib/pq"
)

func Connection(config config.DBconfig) (db *sqlx.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBusername,
		config.DBPassword,
		config.DBName,
	)

	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}

	db.SetMaxIdleConns(int(config.DBPoolConfig.MaxIdleConnection))
	db.SetMaxOpenConns(int(config.DBPoolConfig.MaxOpenConnection))
	db.SetConnMaxLifetime(time.Duration(config.DBPoolConfig.MaxLifetimeConnection))
	db.SetConnMaxIdleTime(time.Duration(config.DBPoolConfig.MaxIdleTimeConnetction))

	return
}
