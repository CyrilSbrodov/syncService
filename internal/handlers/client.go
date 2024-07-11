package handlers

import (
	"encoding/json"
	"errors"
	"github.com/CyrilSbrodov/syncService/internal/model"
	"net/http"
)

// AddClient - ручка добавления нового клиента
func (h *Handler) AddClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var client model.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := h.storage.AddClient(r.Context(), &client); err != nil {
			if errors.Is(err, model.ErrorUserConflict) {
				http.Error(w, "client_name is already exists", http.StatusConflict)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

// UpdateClient - ручка изменения клиента
func (h *Handler) UpdateClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var client model.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := h.storage.UpdateClient(r.Context(), &client); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

// DeleteClient - ручка удаления клиента
func (h *Handler) DeleteClient() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var client model.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := h.storage.DeleteClient(r.Context(), &client); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
