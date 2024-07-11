package storage

import (
	"context"
	"github.com/CyrilSbrodov/syncService/internal/model"
)

// Storage - интерфейс БД
type Storage interface {
	AddClient(ctx context.Context, client *model.Client) error
	UpdateClient(ctx context.Context, client *model.Client) error
	DeleteClient(ctx context.Context, client *model.Client) error
	UpdateAlgorithmStatus(ctx context.Context, as *model.AlgorithmStatus) error
	GetAlgorithmStatus(ctx context.Context) ([]model.AlgorithmStatus, error)
}
