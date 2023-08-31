package database

import (
	"avito_test_case/config"
	"avito_test_case/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx-zerolog"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"strings"
	"time"
)

func NewPostgres(cfg config.Postgres, l logger.Logger) *pgxpool.Pool {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB)

	var db *pgxpool.Pool
	var poolCfg *pgxpool.Config
	var err error
	poolCfg, err = pgxpool.ParseConfig(connStr)
	poolCfg.MaxConns = cfg.PoolMaxOpen
	poolCfg.MinConns = cfg.PoolMaxIdle
	poolCfg.MaxConnLifetime = time.Second * time.Duration(cfg.PoolMaxLifetime)
	switch v := l.(type) {
	case *logger.ZeroLogLogger:
		dbLogLevel, parseErr := tracelog.LogLevelFromString(strings.ToLower(cfg.LogLevel))
		if parseErr != nil {
			l.Fatal("Couldn't parse log level for Postgres", err)
		}
		poolCfg.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   zerolog.NewLogger(*v.L),
			LogLevel: dbLogLevel,
		}
	}
	db, err = pgxpool.NewWithConfig(context.TODO(), poolCfg)
	err = db.Ping(context.TODO())
	if err != nil {
		att := cfg.ReconnectAttempts
		l.Warn("Couldn't connect to Postgres DB")
		for att > 0 && err != nil {
			l.Warn("Attempting Postgres reconnect...")
			db, err = pgxpool.NewWithConfig(context.TODO(), poolCfg)
			err = db.Ping(context.TODO())
			att--
		}
		if err != nil {
			l.Fatal("Couldn't connect to Postgres after %d %s %s. Error: ", cfg.ReconnectAttempts, "attempts", err)
		}
	}

	return db
}
