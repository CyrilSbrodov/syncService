package postgres

import (
	"context"
	"database/sql"
	"github.com/CyrilSbrodov/syncService/cmd/loggers"
	"github.com/CyrilSbrodov/syncService/internal/config"
	"github.com/CyrilSbrodov/syncService/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
	"time"
)

func newMock() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

func TestPGStore_AddClient(t *testing.T) {
	db, mock, err := newMock()
	require.NoError(t, err)
	defer db.Close()

	logger := &loggers.Logger{}
	cfg := &config.Config{}

	store := &PGStore{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}

	client := &model.Client{
		ClientName:  "TestClient",
		Version:     1,
		Image:       "test/image",
		CPU:         "1",
		Memory:      "1Gi",
		Priority:    1.0,
		NeedRestart: false,
		SpawnedAt:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectQuery("INSERT INTO clients").
		WithArgs(client.ClientName, client.Version, client.Image, client.CPU, client.Memory, client.Priority, client.NeedRestart, client.SpawnedAt, client.CreatedAt, client.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("INSERT INTO algorithm_status").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.AddClient(context.Background(), client)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), client.ID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGStore_DeleteClient(t *testing.T) {
	db, mock, err := newMock()
	require.NoError(t, err)
	defer db.Close()

	logger := &loggers.Logger{}
	cfg := &config.Config{}

	store := &PGStore{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}

	client := &model.Client{
		ID: 1,
	}

	mock.ExpectExec("DELETE FROM clients").
		WithArgs(client.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("DELETE FROM algorithm_status").
		WithArgs(client.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.DeleteClient(context.Background(), client)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGStore_UpdateAlgorithmStatus(t *testing.T) {
	db, mock, err := newMock()
	require.NoError(t, err)
	defer db.Close()

	logger := &loggers.Logger{}
	cfg := &config.Config{}

	store := &PGStore{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}

	as := &model.AlgorithmStatus{
		ClientID: 1,
		VWAP:     true,
		TWAP:     false,
		HFT:      true,
	}

	mock.ExpectExec("UPDATE algorithm_status").
		WithArgs(as.VWAP, as.TWAP, as.HFT, as.ClientID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.UpdateAlgorithmStatus(context.Background(), as)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGStore_UpdateClient(t *testing.T) {
	db, mock, err := newMock()
	require.NoError(t, err)
	defer db.Close()

	logger := &loggers.Logger{}
	cfg := &config.Config{}

	store := &PGStore{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}

	client := &model.Client{
		ID:          1,
		ClientName:  "UpdatedClient",
		Version:     2,
		Image:       "updated/image",
		CPU:         "2",
		Memory:      "2Gi",
		Priority:    2.0,
		NeedRestart: true,
		SpawnedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectExec("UPDATE clients").
		WithArgs(client.ClientName, client.Version, client.Image, client.CPU, client.Memory, client.Priority, client.NeedRestart, client.SpawnedAt, sqlmock.AnyArg(), client.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = store.UpdateClient(context.Background(), client)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
