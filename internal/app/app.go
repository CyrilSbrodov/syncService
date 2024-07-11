package app

import (
	"context"
	"errors"
	"github.com/CyrilSbrodov/syncService/cmd/loggers"
	"github.com/CyrilSbrodov/syncService/internal/config"
	"github.com/CyrilSbrodov/syncService/internal/handlers"
	"github.com/CyrilSbrodov/syncService/internal/storage/postgres"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerApp struct {
	cfg    config.Config
	logger *loggers.Logger
	router *mux.Router
}

func NewServerApp() *ServerApp {
	cfg := config.NewConfig()
	router := mux.NewRouter()
	logger := loggers.SetupLogger(cfg.Env)

	return &ServerApp{
		router: router,
		cfg:    *cfg,
		logger: logger,
	}
}

func (a *ServerApp) Run() {
	db, err := postgres.NewPGStore(&a.cfg, a.logger)
	if err != nil {
		a.logger.Error("failed to start pg store", err)
		return
	}
	h := handlers.NewHandler(&a.cfg, a.logger, db)

	h.Register(a.router)

	srv := &http.Server{
		Addr:         a.cfg.Listener.Addr,
		Handler:      a.router,
		ReadTimeout:  a.cfg.Listener.Timeout,
		WriteTimeout: a.cfg.Listener.Timeout,
		IdleTimeout:  a.cfg.Listener.IdleTimeout,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			a.logger.Error("server not started", err, "server")
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		a.logger.Error("server", "failed to shutting down gracefully", err)
		return
	}
	a.logger.Info("shutting down", slog.String("server", a.cfg.Listener.Addr))
	os.Exit(0)
}
