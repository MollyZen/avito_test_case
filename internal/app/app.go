package app

import (
	"avito_test_case/config"
	"avito_test_case/internal/datastruct"
	"avito_test_case/internal/repository"
	"avito_test_case/pkg/database"
	"avito_test_case/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l := logger.NewZeroLog(cfg.Log.Level)
	db := database.NewPostgres(cfg.Postgres)
	segRep := repository.NewDatabaseSegmentRepository(db, &l)
	_, err := segRep.Create(context.TODO(), datastruct.Segment{
		Name: "TEST_SEG",
	})
	if err != nil {
		l.Fatal("Error adding new segment: ", err)
	}

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("main - Run - signal: " + s.String())
	}
}
