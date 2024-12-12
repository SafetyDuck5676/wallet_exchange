--- 
 
### 9. README.md 
 
Создаем README.md с описанием проекта: 
 
markdown 
# gw-exchanger 
 
gRPC-сервис для получения курсов валют. 
 
## Функционал 
- Получение всех курсов валют. 
- Получение курса обмена между двумя валютами. 
 
## Структура проекта  
gw-exchanger/ 
├── cmd/ 
│   └── main.go 
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
├── proto/ 
│   └── exchange.proto 
├── config.env 
├── Makefile 
└── Dockerfile 
 
 
## Запуск 
1. Установите зависимости: 
   ```bash 
   go mod tidy 
    
2. Соберите и запустите: 
   bash 
   make build 
   make run 
    
3. Используйте Docker: 
   bash 
   make docker-build 
   make docker-run 
    
 
## Переменные окружения 
- DATABASE_URL — строка подключения к PostgreSQL. 
- GRPC_PORT — порт для gRPC-сервера. 
 
## gRPC API 
- GetExchangeRates — возвращает курсы всех валют. 
- GetExchangeRateForCurrency — возвращает курс для двух валют. 
 
 
 
--- 
 
### 10. Логирование (`pkg/utils.go`) 
 
Добавляем логирование для всего сервиса: 
 
```go 
package utils 
 
import ( 
 "log" 
 "os" 
) 
 
var ( 
 InfoLogger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile) 
 ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile) 
) 
 
func LogInfo(message string) { 
 InfoLogger.Println(message) 
} 
 
func LogError(err error) { 
 ErrorLogger.Println(err) 
}  
 
Используем логирование в проекте. Например, в cmd/main.go: 
 
go 
import "gw-exchanger/pkg/utils" 
 
func main() { 
 utils.LogInfo("Starting gw-exchanger service...") 
 // Остальной код... 
}  
 
--- 
 
### 11. Инициализация базы данных 
 
Скрипт для создания таблицы в PostgreSQL (`db/init.sql`): 
 
sql 
CREATE TABLE exchange_rates ( 
    id SERIAL PRIMARY KEY, 
    from_currency VARCHAR(3) NOT NULL, 
    to_currency VARCHAR(3) NOT NULL, 
    rate NUMERIC(10, 5) NOT NULL 
); 
 
INSERT INTO exchange_rates (from_currency, to_currency, rate) 
VALUES 
('USD', 'EUR', 0.85), 
('EUR', 'USD', 1.18), 
('USD', 'RUB', 70.00), 
('RUB', 'USD', 0.014), 
('EUR', 'RUB', 80.00), 
('RUB', 'EUR', 0.0125);  
 
---