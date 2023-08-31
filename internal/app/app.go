package app

import (
	"avito_test_case/config"
	_ "avito_test_case/docs"
	"avito_test_case/internal/controller/http/v1"
	"avito_test_case/internal/repository"
	"avito_test_case/internal/service"
	"avito_test_case/pkg/database"
	"avito_test_case/pkg/httpserver"
	"avito_test_case/pkg/logger"
	"context"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	asService := service.NewPostgresAssignmentService(db, usRep, segRep, asRep, hisRep, l)
	segService := service.NewPostgresSegmentService(db, segRep, asRep, hisRep, l)

	//http
	mux := v1.NewRouter(l)
	v1.NewSegmentController(mux, segService, l)
	v1.NewAssignmentController(mux, asService, l)
	v1.NewUserController(mux, asService, l)
	mux.Mount("/swagger", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	mux.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
	httpServer := httpserver.New(mux, cfg.HTTP)

	//delete expired assignments every hour
	ticker := time.NewTicker(time.Duration(cfg.App.AssignmentCleaningInterval) * time.Second)
	cleaner := make(chan interface{})
	go func() {
		for range ticker.C {
			l.Info("Cleaning expired assignments...")
			err := asRep.DeleteExpired(context.TODO())
			if err != nil {
				cleaner <- err
			}
		}
	}()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	l.Info("The app is running")

	select {
	case s := <-interrupt:
		l.Info("main - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err := <-cleaner:
		l.Error(fmt.Errorf("app - Run - Cleaner: %w", err))
	}

	// Shutdown

	if err := httpServer.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
	ticker.Stop()
}
