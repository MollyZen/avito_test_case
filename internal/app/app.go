package app

import (
	"avito_test_case/config"
	"avito_test_case/internal/repository"
	"avito_test_case/internal/service"
	"avito_test_case/pkg/database"
	"avito_test_case/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l := logger.NewZeroLog(cfg.Log.Level)

	db := database.NewPostgres(cfg.Postgres, l)

	//repos
	usRep := repository.NewPostgresUserRepository(db, l)
	segRep := repository.NewPostgresSegmentRepository(db, l)
	asRep := repository.NewPostgresAssignmentRepository(db, l)
	hisRep := repository.NewPostgresHistoryRepository(db, l)

	//services
	service.NewAssignmentUseCase(usRep, segRep, asRep, hisRep, l)

	//http

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("main - Run - signal: " + s.String())
	}
}
