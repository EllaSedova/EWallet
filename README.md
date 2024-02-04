# EWallet
Приложение представляет собой HTTP сервер, реализующий систему обработки транзакций платежной системы. Приложение реализовано в виде REST API, которое предоставляет 5 методов для работы с кошельками и их транзакциями.
Для безопасной передачи информации между клиентом и сервером используется JSON Web Token.

## Установка
1. Установите Go 1.21
2. Установите PostgreSQL и создайте базу данных
3. Склонируйте репозиторий EWallet:  
https://github.com/EllaSedova/EWallet.git
4. Установите зависимости проекта:
```
go mod download
```
5. Внесите актуальные данные в файл .env:
```
db_name = name
db_pass = password
db_user = user
db_type = type
db_host = localhost
db_port = 5434
token_password = thisIsTheJwtSecretPassword
```
6. Запустите приложение:
```
go run main.go
```

## Использование
После установки и запуска приложения, вы можете обращаться к нему с помощью HTTP запросов.
### Создание кошелька
```
POST /api/v1/wallet
```
Параметры запроса: Отсутствуют  
Пример запроса:
```
POST /api/v1/wallet
```
Пример ответа:
```json
{
    "description": "Кошелёк создан",
    "status": 200,
    "wallet": {
        "id": "e68360f4-bb53-43dd-a24d-2a38f28e1f80",
        "balance": 100,
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJXYWxsZXRJZCI6ImU2ODM2MGY0LWJiNTMtNDNkZC1hMjRkLTJhMzhmMjhlMWY4MCJ9.sgJU497EWOyz3xoNu1BvWux2hja-JHNkEHqmjjBoRkw"
    }
}
```
status code: 200 OK

### Перевод средств
```
POST /api/v1/wallet/{walletId}/send
```
Параметры запроса:
- walletId: строковый ID кошелька, указанный в пути запроса
- JSON-объект в теле запроса с параметрами:
  - to: ID кошелька, куда нужно перевести деньги
  - amount: сумма перевода

Пример запроса:
```
POST /api/v1/wallet/e68360f4-bb53-43dd-a24d-2a38f28e1f80/send
```
```json
Content-Type: application/json
{
  "to": "3fa2d6f8-120a-4400-bd85-79107b5e179d",
  "amount": 50.0
}
```
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJXYWxsZXRJZCI6ImU2ODM2MGY0LWJiNTMtNDNkZC1hMjRkLTJhMzhmMjhlMWY4MCJ9.sgJU497EWOyz3xoNu1BvWux2hja-JHNkEHqmjjBoRkw
```
Пример ответа:
```json
{
    "description": "Перевод успешно проведен",
    "status": 200,
}
```
status code: 200 OK

### Получение истории транзакций
```
GET /api/v1/wallet/{walletId}/history
```
Параметры запроса:
- walletId: Строковый ID кошелька, указанный в пути запроса

Пример запроса:
```
GET /api/v1/wallet/25197d19-488a-4ee5-a0d9-f333f9623a0f/history
```
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJXYWxsZXRJZCI6IjI1MTk3ZDE5LTQ4OGEtNGVlNS1hMGQ5LWYzMzNmOTYyM2EwZiJ9.S1SOe1-elNEQm-Vq3YlGM5GGqX9QvIOxS0C0LmLyVpc
```

Пример ответа:
```json
{
    "description": "История транзакций получена",
    "in": [
        {
            "ID": 6,
            "from": "3fa2d6f8-120a-4400-bd85-79107b5e179d",
            "to": "25197d19-488a-4ee5-a0d9-f333f9623a0f",
            "amount": 2,
            "time": "2024-02-04T01:24:10.564222+03:00"
        }
    ],
    "out": [
        {
            "ID": 4,
            "from": "25197d19-488a-4ee5-a0d9-f333f9623a0f",
            "to": "3fa2d6f8-120a-4400-bd85-79107b5e179d",
            "amount": 3,
            "time": "2024-02-03T21:06:41.853178+03:00"
        }
    ],
    "status": 200
}
```
status code: 200 OK

### Получение состояния кошелька
```
GET /api/v1/wallet/{walletId}
```
Параметры запроса:
- walletId: Строковый ID кошелька, указанный в пути запроса

Пример запроса:
```
GET /api/v1/wallet/25197d19-488a-4ee5-a0d9-f333f9623a0f
```
```
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJXYWxsZXRJZCI6IjI1MTk3ZDE5LTQ4OGEtNGVlNS1hMGQ5LWYzMzNmOTYyM2EwZiJ9.S1SOe1-elNEQm-Vq3YlGM5GGqX9QvIOxS0C0LmLyVpc
```
Пример ответа:
```json
{
    "data": {
        "id": "25197d19-488a-4ee5-a0d9-f333f9623a0f",
        "balance": 7000,
        "token": ""
    },
    "description": "Данные кошелька получены",
    "status": 200
}
```
status code: 200 OK

### Аутентификация пользователя (необходима для того, чтобы перегенерировать JWT)
```
POST /api/v1/wallet/login 
```
Параметры запроса:
- JSON-объект в теле запроса с параметрами:
  - id: ID кошелька, для которого нужно перегенерировать JWT

Пример запроса:
```
POST /api/v1/wallet/login
```
```json
Content-Type: application/json
{
  "id": "56f4b6dd-9071-433d-94df-77a0a36c71bc"
}
```
Пример ответа:
```json
{
    "description": "JWT токен обновлён",
    "status": 200,
    "wallet": {
        "id": "56f4b6dd-9071-433d-94df-77a0a36c71bc",
        "balance": 88,
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJXYWxsZXRJZCI6IjU2ZjRiNmRkLTkwNzEtNDMzZC05NGRmLTc3YTBhMzZjNzFiYyJ9.s4i0p_rxgxGlmzieM1jkILBLwDbjBL4UM14vJ7z7ceE"
    }
}
```
status code: 200 OK
## База данных
Данные кошельков и транзакций сохраняются в базе данных.  
На основании учётных данных из файла .env производится подключение к БД. При первом запуске автоматически создаются таблицы:

**"wallets"**:
- Атрибуты:
  - _id_ (строковый ID кошелька, генерируется сервером) PK
  - _balance_ (дробное число, баланс кошелька) not null

**"transactions"**:
- Атрибуты:
  - _id_ (уникальный идентификатор транзакции) PK
  - _from_ (строковый ID исходящего кошелька) not null
  - _to_ (строковый ID входящего кошелька) not null
  - _amount_ (дробное число, сумма перевода) not null
  - _time_ (дата и время перевода) default now()
