package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CyrilSbrodov/syncService/internal/model"
	"github.com/stretchr/testify/assert"
)

type mockStorage struct {
	addClient             func(ctx context.Context, client *model.Client) error
	updateClient          func(ctx context.Context, client *model.Client) error
	deleteClient          func(ctx context.Context, client *model.Client) error
	updateAlgorithmStatus func(ctx context.Context, a *model.AlgorithmStatus) error
	getAlgorithmStatus    func(ctx context.Context) ([]model.AlgorithmStatus, error)
}

func (m *mockStorage) AddClient(ctx context.Context, client *model.Client) error {
	return m.addClient(ctx, client)
}

func (m *mockStorage) UpdateClient(ctx context.Context, client *model.Client) error {
	return m.updateClient(ctx, client)
}

func (m *mockStorage) DeleteClient(ctx context.Context, client *model.Client) error {
	return m.deleteClient(ctx, client)
}

func (m *mockStorage) UpdateAlgorithmStatus(ctx context.Context, a *model.AlgorithmStatus) error {
	return m.updateAlgorithmStatus(ctx, a)
}

func (m *mockStorage) GetAlgorithmStatus(ctx context.Context) ([]model.AlgorithmStatus, error) {
	return m.getAlgorithmStatus(ctx)
}

func TestAddClient(t *testing.T) {
	tests := []struct {
		name           string
		inputBody      *model.Client
		storageError   error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "200",
			inputBody:      &model.Client{ClientName: "Client"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "400",
			inputBody:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body\n",
		},
		{
			name:           "409",
			inputBody:      &model.Client{ClientName: "Client1"},
			storageError:   model.ErrorClientConflict,
			expectedStatus: http.StatusConflict,
			expectedBody:   "client_name is already exists\n",
		},
		{
			name:           "500",
			inputBody:      &model.Client{ClientName: "Client2"},
			storageError:   errors.New("error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &mockStorage{
				addClient: func(ctx context.Context, client *model.Client) error {
					return tt.storageError
				},
			}
			handler := &Handler{storage: storage}
			rr := httptest.NewRecorder()
			var req *http.Request
			if tt.inputBody != nil {
				body, err := json.Marshal(tt.inputBody)
				assert.NoError(t, err)
				req = httptest.NewRequest(http.MethodPost, "/client", bytes.NewBuffer(body))
			} else {
				req = httptest.NewRequest(http.MethodPost, "/client", nil)
			}

			handler.AddClient()(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestUpdateClient(t *testing.T) {
	tests := []struct {
		name           string
		inputBody      *model.Client
		storageError   error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "200",
			inputBody:      &model.Client{ClientName: "Client"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "400",
			inputBody:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body\n",
		},
		{
			name:           "500",
			inputBody:      &model.Client{ClientName: "Client1"},
			storageError:   errors.New("error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &mockStorage{
				updateClient: func(ctx context.Context, client *model.Client) error {
					return tt.storageError
				},
			}
			handler := &Handler{storage: storage}
			rr := httptest.NewRecorder()
			var req *http.Request
			if tt.inputBody != nil {
				body, err := json.Marshal(tt.inputBody)
				assert.NoError(t, err)
				req = httptest.NewRequest(http.MethodPut, "/client", bytes.NewBuffer(body))
			} else {
				req = httptest.NewRequest(http.MethodPut, "/client", nil)
			}

			handler.UpdateClient()(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestDeleteClient(t *testing.T) {
	tests := []struct {
		name           string
		inputBody      *model.Client
		storageError   error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "200",
			inputBody:      &model.Client{ClientName: "Client"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "400",
			inputBody:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body\n",
		},
		{
			name:           "500",
			inputBody:      &model.Client{ClientName: "Client1"},
			storageError:   errors.New("error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &mockStorage{
				deleteClient: func(ctx context.Context, client *model.Client) error {
					return tt.storageError
				},
			}
			handler := &Handler{storage: storage}
			rr := httptest.NewRecorder()
			var req *http.Request
			if tt.inputBody != nil {
				body, err := json.Marshal(tt.inputBody)
				assert.NoError(t, err)
				req = httptest.NewRequest(http.MethodDelete, "/client", bytes.NewBuffer(body))
			} else {
				req = httptest.NewRequest(http.MethodDelete, "/client", nil)
			}

			handler.DeleteClient()(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}
