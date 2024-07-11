package syncer

import (
	"context"
	"fmt"
	"github.com/CyrilSbrodov/syncService/cmd/loggers"
	"github.com/CyrilSbrodov/syncService/internal/config"
	"github.com/CyrilSbrodov/syncService/internal/deployer"
	"github.com/CyrilSbrodov/syncService/internal/storage"
	"time"
)

// Syncer - структура синкера, что своими методами проверяет запущенные pods
type Syncer struct {
	store    storage.Storage
	deployer deployer.Deployer
	logger   *loggers.Logger
	cfg      config.Config
}

// NewSyncer - конструктор синкера
func NewSyncer(d deployer.Deployer, store storage.Storage, logger *loggers.Logger, cfg config.Config) *Syncer {
	return &Syncer{
		store:    store,
		deployer: d,
		logger:   logger,
		cfg:      cfg,
	}
}

// Start - функция запуска синкера с таймером на 5 минут
func (s *Syncer) Start() {
	ticker := time.NewTicker(s.cfg.SyncTimeout)
	for {
		select {
		case <-ticker.C:
			s.syncAlgorithms()
		}
	}
}

// syncAlgorithms - функция синхронизации алгоритмов с базой данных
func (s *Syncer) syncAlgorithms() {
	algorithms, err := s.store.GetAlgorithmStatus(context.Background())
	if err != nil {
		s.logger.Error("Error fetching clients:", err)
		return
	}

	list, _ := s.deployer.GetPodList()
	fmt.Println(list)

	if algorithms == nil {
		return
	}
	for _, a := range algorithms {
		if a.VWAP {
			err := s.deployer.CreatePod(fmt.Sprintf("vmap-%d", a.AlgorithmID))
			if err != nil {
				s.logger.Error("Error creating pod:", err)
				return
			}
		} else {
			err := s.deployer.DeletePod(fmt.Sprintf("vmap-%d", a.AlgorithmID))
			if err != nil {
				s.logger.Error("Error delete pod:", err)
				return
			}
		}
		if a.HFT {
			err := s.deployer.CreatePod(fmt.Sprintf("hft-%d", a.AlgorithmID))
			if err != nil {
				s.logger.Error("Error creating pod:", err)
				return
			}
		} else {
			err := s.deployer.DeletePod(fmt.Sprintf("hft-%d", a.AlgorithmID))
			if err != nil {
				s.logger.Error("Error delete pod:", err)
				return
			}
		}
		if a.TWAP {
			err := s.deployer.CreatePod(fmt.Sprintf("twap-%d", a.AlgorithmID))
			if err != nil {
				s.logger.Error("Error creating pod:", err)
				return
			}
		} else {
			err := s.deployer.DeletePod(fmt.Sprintf("twap-%d", a.AlgorithmID))
			if err != nil {
				s.logger.Error("Error delete pod:", err)
				return
			}
		}
	}
}
