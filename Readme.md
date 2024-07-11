#### Оглавление:
____
0. [Сервис синхронизации алгоритмов](#сервис-синхронизации-алгоритмов).
1. [ЗАВИСИМОСТИ](#зависимости).
2. [ЗАПУСК/СБОРКА](#запусксборка).
   2.1 [Конфигурация](#конфигурация).
   2.2 [Запуск сервера](#запуск-сервера).
3. [О сервисе](#о-сервисе).
____

# Сервис синхронизации алгоритмов.

Сервис позволяет создавать, удалять, изменять клиентов и их алгоритмы. В автоматическом режиме, раз в 5 минут сервис проверяет базу данных, если в базе он находит алгоритмы со статусом true, то запускает pod. Если статус pod false, то сервис удаляет pod, если такой ранее был создан, или ничего не делает.

Структура [клиента и алгоритмов](https://github.com/CyrilSbrodov/syncService/blob/main/internal/model/models.go):
```GO
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
```

Структура сервиса следующая:
1) Сервер - обработка полученных данных и отправка их в БД Postgres.
2) Синкер - проверка состояния алгоритмов (создание или удаление pods).
3) БД - прием получаемых данных.
____
# ЗАВИСИМОСТИ.

Используется язык go версии 1.18. Используемые библиотеки:
- github.com/gorilla/mux v1.8.1
- github.com/ilyakaznacheev/cleanenv v1.5.0
- github.com/lib/pq v1.10.9
- github.com/stretchr/testify v1.9.0
- k8s.io/api v0.30.2
- k8s.io/apimachinery v0.30.2
- k8s.io/client-go v0.30.2
- POSTGRESQL latest
____

# ЗАПУСК/СБОРКА.

## Конфигурация.

Предусмотрены различные конфигурации:
1) флаги
2) облачные переменные
3) конфигурационный файл

## Запуск сервера.

Есть несколько способов запуска:
- Если у Вас установлен make, то можно использовать Makefile. Необходимо в консоле прописать:
```
make run
```

- Можно запустить сервер из пакета [cmd](https://github.com/CyrilSbrodov/syncService/blob/main/cmd/main.go)
```
cd cmd
go run main.go
```
- Если у Вас не установлена База PostgreSQL, то можно запустить файл docker-compose командой:
```
docker compose up -d
либо
docker-compose up -d
```

# О сервсие.
Структура приложения позволяет нативно вносить корректировки:

[Структура сервера](https://github.com/CyrilSbrodov/syncService/blob/main/internal/app/app.go):
```GO
// ServerApp - структура сервера
type ServerApp struct {
    cfg    config.Config
    logger *loggers.Logger
    router *mux.Router
}
```

Немного о сервере:

Сервер получает данные по следующим эндпоинтам:
[http](https://github.com/CyrilSbrodov/syncService/blob/main/internal/handlers/handler.go):
```GO
func (h *Handler) Register(r *mux.Router) {
    r.HandleFunc("/api/client", h.AddClient()).Methods("POST")
    r.HandleFunc("/api/client", h.UpdateClient()).Methods("PUT")
    r.HandleFunc("/api/client", h.DeleteClient()).Methods("DELETE")
    r.HandleFunc("/api/algorithms", h.UpdateAlgorithmStatus()).Methods("POST")
}
```