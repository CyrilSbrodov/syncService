package handlers

import (
	"github.com/CyrilSbrodov/syncService/cmd/loggers"
	"github.com/CyrilSbrodov/syncService/internal/config"
	"github.com/CyrilSbrodov/syncService/internal/storage"
	"github.com/gorilla/mux"
)

type Handlers interface {
	Register(router *mux.Router)
}

type Handler struct {
	cfg     *config.Config
	logger  *loggers.Logger
	storage storage.Storage
}

func NewHandler(cfg *config.Config, logger *loggers.Logger, storage storage.Storage) *Handler {
	return &Handler{
		cfg:     cfg,
		logger:  logger,
		storage: storage,
	}
}

func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/api/register", h.SignUp()).Methods("POST")
	r.HandleFunc("/api/login", h.SignIn()).Methods("POST")
	secure := r.PathPrefix("/auth").Subrouter()
	secure.Use(h.userIdentity)
	secure.HandleFunc("/api/task", h.GetAll()).Methods("GET")
	secure.HandleFunc("/api/task", h.NewList()).Methods("POST")
}
