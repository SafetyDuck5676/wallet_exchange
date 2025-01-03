# Сервис кошелька

## Обзор
Сервис кошелька – это микросервис, отвечающий за управление кошельками пользователей, обработку транзакций и ведение записей о балансе. Он является частью проекта Docker Exchanger.

## Возможности
- Создание и управление кошельками пользователей.
- Обработка депозитов и снятий.
- Обмен денег между различными валютами.
- Обеспечение согласованности и целостности данных.

## Требования
- Docker.
- Docker Compose.

## Использование
- Доступ к API сервиса кошелька осуществляется по адресу: `http://localhost:8080/api/v1`.

## API Эндпоинты
- `POST /register` - Создание новой учетной записи пользователя с кошельками в валютах RUB, USD и EUR.
- `POST /login` - Вход пользователя с использованием имени пользователя и пароля. Возвращает JWT-токен для авторизации в API.
- `POST /balance` - Требует JWT-токен, возвращает средства на кошельках пользователя.
- `POST /wallet/deposit` - Вносит деньги в кошелек с указанной валютой.
- `POST /wallet/withdraw` - Снимает деньги с кошелька с указанной валютой.
- `GET /rates` - Возвращает текущие курсы обмена от сервера обменника.
- `POST /rate` - Возвращает курс обмена одной валюты на другую.
- `POST /exchange` - Снимает деньги с одного кошелька и зачисляет эквивалентную сумму на кошелек с другой валютой.

## Детальное описание
Эндпоинт `register` API создает нового пользователя, три записи в таблице кошельков и три записи в таблице балансов, ссылаясь на таблицу валют для соответствующей валюты кошелька.

### Вход
Если вход выполнен успешно, ID пользователя и имя пользователя шифруются в JWT-токене. Этот токен требуется для всех последующих вызовов API.

### Баланс
Эндпоинт `balance` выполняет простой запрос к таблице, хранящей данные о пользователе, который идентифицируется с помощью JWT-токена.

### Депозит, снятие и обмен
При выполнении операций депозитов, снятия и обмена ожидается, что данные будут переданы в виде числа с плавающей точкой (float). Однако деньги хранятся в базе данных как целые числа (integer) для обеспечения точности. При обмене денег все значения с плавающей точкой преобразуются в целые числа с использованием коэффициента преобразования. Аналогичным образом происходит преобразование обратно в числа с плавающей точкой.

Этот подход позволяет избежать ошибок округления, а также обеспечивает более быструю обработку операций с целыми числами по сравнению с числами с плавающей точкой.

### Архитектура сервиса
Сервис разделен на три части: обработчик (handler), сервис (service) и репозиторий (repository).
- **Обработчики** вызываются HTTP-запросами через маршруты, определенные в `main.go`. Используются для проверки токенов, обработки запросов и отправки ответов пользователю API.
- **Сервисы** используются как промежуточное звено для соединения обработчиков API с функциональностью репозитория.
- **Функции репозитория** реализуют основную логику сервиса, включая SQL-запросы, создание токенов, регистрацию новых пользователей и сбор данных с сервера обменника.
