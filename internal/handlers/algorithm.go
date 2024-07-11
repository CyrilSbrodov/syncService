package handlers

import (
	"encoding/json"
	"github.com/CyrilSbrodov/syncService/internal/model"
	"net/http"
)

// UpdateAlgorithmStatus - ручка обновления статусов алгоритмов
func (h *Handler) UpdateAlgorithmStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var as model.AlgorithmStatus
		if err := json.NewDecoder(r.Body).Decode(&as); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := h.storage.UpdateAlgorithmStatus(r.Context(), &as); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
