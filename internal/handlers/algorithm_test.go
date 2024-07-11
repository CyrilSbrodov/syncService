package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/CyrilSbrodov/syncService/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_UpdateAlgorithmStatus(t *testing.T) {
	tests := []struct {
		name               string
		inputBody          interface{}
		mockUpdateFunc     func(ctx context.Context, as *model.AlgorithmStatus) error
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "ok",
			inputBody: model.AlgorithmStatus{
				ClientID: 1,
				VWAP:     true,
				TWAP:     false,
				HFT:      true,
			},
			mockUpdateFunc: func(ctx context.Context, as *model.AlgorithmStatus) error {
				return nil
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "",
		},
		{
			name:               "400",
			inputBody:          "invalid body",
			mockUpdateFunc:     nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "invalid request body\n",
		},
		{
			name: "500",
			inputBody: model.AlgorithmStatus{
				ClientID: 1,
				VWAP:     true,
				TWAP:     false,
				HFT:      true,
			},
			mockUpdateFunc: func(ctx context.Context, as *model.AlgorithmStatus) error {
				return fmt.Errorf("error from db")
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "internal server error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &mockStorage{
				updateAlgorithmStatus: tt.mockUpdateFunc,
			}
			handler := &Handler{
				storage: mockStorage,
			}
			jsonBody, err := json.Marshal(tt.inputBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("PUT", "/api/algorithms", bytes.NewBuffer(jsonBody))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.UpdateAlgorithmStatus()(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}
