package database

import (
	"avito_test_case/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func NewPostgres(cfg config.Postgres) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB)

	var db *sql.DB
	var err error
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		att := cfg.ReconnectAttempts
		for att > 0 && err != nil {
			db, err = sql.Open("postgres", connStr)
			log.Print("Attempting reconnect...")
			att--
		}
		if err != nil {
			panic(err)
		}
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Second * time.Duration(cfg.PoolMaxLifetime))
	db.SetMaxOpenConns(cfg.PoolMaxOpen)
	db.SetMaxIdleConns(cfg.PoolMaxIdle)

	return db
}
