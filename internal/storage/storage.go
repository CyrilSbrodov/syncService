package storage

// Storage - интерфейс БД
type Storage interface {
	AddClient()
	UpdateClient()
	DeleteClient()
	UpdateAlgorithmStatus()
}
