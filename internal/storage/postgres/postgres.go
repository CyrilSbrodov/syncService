package postgres

import (
	"context"
	"database/sql"
	"github.com/CyrilSbrodov/syncService/cmd/loggers"
	"github.com/CyrilSbrodov/syncService/internal/config"
	"github.com/CyrilSbrodov/syncService/internal/model"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"time"
)

// PGStore - структура БД
type PGStore struct {
	cfg    *config.Config
	logger *loggers.Logger
	db     *sql.DB
}

// NewPGStore - конструктор БД
func NewPGStore(cfg *config.Config, logger *loggers.Logger) (*PGStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", cfg.DBPath)
	if err != nil {
		logger.Error("failed to init db", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		logger.Error("failed to connect to db", err)
		return nil, err
	}
	if err := createTable(ctx, db, logger); err != nil {
		logger.Error("failed to create table", err)
		return nil, err
	}
	return &PGStore{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}, nil
}

// createTable - функция создания новых таблиц в БД.
func createTable(ctx context.Context, db *sql.DB, logger *loggers.Logger) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		logger.Error("failed to begin transaction", err)
		return err
	}

	defer tx.Rollback()

	//создание таблиц
	tables := []string{
		`CREATE TABLE IF NOT EXISTS clients (
			id SERIAL PRIMARY KEY,
			client_name VARCHAR(100) NOT NULL UNIQUE,
			version INT,
			image VARCHAR(255),
			cpu VARCHAR(50),
    		memory VARCHAR(50),
    		priority FLOAT,
    		need_restart BOOLEAN DEFAULT FALSE,
    		spawned_at TIMESTAMP,
    		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS algorithm_status (
			id SERIAL PRIMARY KEY,
    		client_id INT REFERENCES clients(id),
    		vwap BOOLEAN DEFAULT FALSE,
    		twap BOOLEAN DEFAULT FALSE,
    		hft BOOLEAN DEFAULT FALSE
		)`,
	}

	for _, table := range tables {
		if _, err = tx.ExecContext(ctx, table); err != nil {
			logger.Error("Unable to create table", err)
			return err
		}
	}

	return tx.Commit()
}

// AddClient - добаление клиента в БД и дефолтные значения алгоритмов
func (p *PGStore) AddClient(ctx context.Context, c *model.Client) error {
	q := `INSERT INTO clients (client_name, version, image, cpu, memory, priority, need_restart, spawned_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	err := p.db.QueryRowContext(ctx, q, c.ClientName, c.Version, c.Image, c.CPU, c.Memory, c.Priority,
		c.NeedRestart, c.SpawnedAt, c.CreatedAt, c.UpdatedAt).Scan(&c.ID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			p.logger.Error("client_name already exists", err)
			return model.ErrorClientConflict
		}
		p.logger.Error("Failure to insert client into table", err)
		return err
	}

	q = `INSERT INTO algorithm_status (client_id, vwap, twap, hft) VALUES ($1, default, default, default)`
	if _, err := p.db.ExecContext(ctx, q, c.ID); err != nil {
		p.logger.Error("Failure to insert algorithm status into table", err)
		return err
	}
	return nil
}

// UpdateClient - обновление клиента в БД
func (p *PGStore) UpdateClient(ctx context.Context, client *model.Client) error {
	q := `UPDATE clients SET client_name=$1, version=$2, image=$3, cpu=$4, memory=$5, priority=$6, need_restart=$7, spawned_at=$8, updated_at=$9 
			WHERE id=$10`
	_, err := p.db.ExecContext(ctx, q, client.ClientName, client.Version, client.Image, client.CPU, client.Memory, client.Priority,
		client.NeedRestart, client.SpawnedAt, time.Now(), client.ID)
	if err != nil {
		p.logger.Error("Failure to update client in table", err)
		return err
	}
	return nil
}

// DeleteClient - удаление клиента и алгоритмов из БД
func (p *PGStore) DeleteClient(ctx context.Context, client *model.Client) error {
	q := `DELETE FROM clients WHERE id=$1`
	_, err := p.db.ExecContext(ctx, q, client.ID)
	if err != nil {
		p.logger.Error("Failure to delete client from table", err)
		return err
	}
	q = `DELETE FROM algorithm_status WHERE client_id=$1`
	_, err = p.db.ExecContext(ctx, q, client.ID)
	if err != nil {
		p.logger.Error("Failure to delete algorithm from table", err)
		return err
	}
	return nil
}

// UpdateAlgorithmStatus - обновление статусов алноритмов в БД
func (p *PGStore) UpdateAlgorithmStatus(ctx context.Context, as *model.AlgorithmStatus) error {
	q := `UPDATE algorithm_status SET vwap=$1, twap=$2, hft=$3 WHERE client_id=$4`
	_, err := p.db.ExecContext(ctx, q, as.VWAP, as.TWAP, as.HFT, as.ClientID)
	if err != nil {
		p.logger.Error("Failure to update algorithm status in table", err)
		return err
	}
	return nil
}

func (p *PGStore) GetAlgorithmStatus(ctx context.Context) ([]model.AlgorithmStatus, error) {
	q := `SELECT id, client_id, vwap, twap, hft FROM algorithm_status`
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		p.logger.Error("Failure to select algorithms from table", err)
		return nil, err
	}
	var algorithms []model.AlgorithmStatus
	for rows.Next() {
		var a model.AlgorithmStatus
		if err := rows.Scan(&a.AlgorithmID, &a.ClientID, &a.VWAP, &a.TWAP, &a.HFT); err != nil {
			p.logger.Error("failed to scan algorithms from data")
			return nil, err
		}
		algorithms = append(algorithms, a)
	}
	return algorithms, nil
}
