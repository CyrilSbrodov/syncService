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
	r.HandleFunc("/api/client", h.AddClient()).Methods("POST")
	r.HandleFunc("/api/client", h.UpdateClient()).Methods("PUT")
	r.HandleFunc("/api/client/{id}", h.DeleteClient()).Methods("DELETE")
	r.HandleFunc("/api/algorithms", h.UpdateAlgorithmStatus()).Methods("POST")
}
