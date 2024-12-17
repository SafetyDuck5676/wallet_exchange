# gw-exchanger 
 
gRPC-сервис для получения курсов валют. 
 
## Функционал 
- Получение всех курсов валют. 
- Получение курса обмена между двумя валютами. 
 
## Структура проекта  

```
gw-exchanger/ 
├── cmd/ 
│   └── main.go
├── db/
│   └── init.sql
├── internal/ 
│   ├── storages/ 
│   │   ├── storage.go 
│   │   ├── model.go 
│   │   └── postgres/ 
│   │       ├── connector.go 
│   │       └── methods.go 
│   ├── config/ 
│   │   ├── config.go 
│   │   └── defaults.go 
│   └── grpc/ 
│       ├── server.go 
│       └── handlers.go
│       pb/
│       ├── exchange.pb.go
│       └── exchange_grpc.pb.go 
├── pkg/ 
│   └── utils.go 
├── config.env 
├── Makefile
├── go.mod
└── Dockerfile 
 ```
 
## Запуск 
Нужно поднять docker контейнер с проектом внутри
```
docker compose up
```
## Переменные окружения 
- DATABASE_URL — строка подключения к PostgreSQL. 
- GRPC_PORT — порт для gRPC-сервера. 
 
## gRPC API 
- GetExchangeRates — возвращает курсы всех валют. 
- GetExchangeRateForCurrency — возвращает курс для двух валют. 
