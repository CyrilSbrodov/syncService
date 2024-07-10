package model

import "time"

// Client - структура клиента.
type Client struct {
	ID          int64     `json:"id"`
	ClientName  string    `json:"client_name"`
	Version     int       `json:"version"`
	Image       string    `json:"image"`
	CPU         string    `json:"cpu"`
	Memory      string    `json:"memory"`
	Priority    float64   `json:"priority"`
	NeedRestart bool      `json:"needRestart"`
	SpawnedAt   time.Time `json:"spawned_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AlgorithmStatus - структура алгоритмов
type AlgorithmStatus struct {
	AlgorithmID int64 `json:"algorithm_id"`
	ClientID    int64 `json:"client_id"`
	VWAP        bool  `json:"vwap"`
	TWAP        bool  `json:"twap"`
	HFT         bool  `json:"hft"`
}
