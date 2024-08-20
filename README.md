# Документация к API

## 1. Обзор

Этот API предназначен для работы с базой данных Tarantool, обработки пользовательской аутентификации, управления сессиями и хранения данных в формате ключ-значение.<br>
Для аутентификации используется JWT-токен. У токена есть время действия (выдается на 3 минуты), по истечению сессия удаляется и необходимо заново авторизоваться. Регистрация не предусмотрена, существует только один пользователь по умолчанию (`username: admin, password: presale`).<br>
Сервис предоставляет API для входа, записи и чтения данных, соответственно по `/api/login`, `/api/write`, `/api/read`.<br>
Все сессии, пользователи и KV хранятся в Tarantool.

## 2. Предварительные требования

Перед запуском убедитесь, что у вас установлен Docker.

## 3. Быстрый запуск

1. Клонируйте репозиторий:

```bash
git clone https://github.com/tarbeevms/internship-vk
```

2. Войдите в скопированную директорию в папку ./build:

```bash
cd internship-vk/build
```

3. Запустите Docker Compose:

```bash
docker compose up 
```

После запуска docker compose запустятся два контейнера: tarantool и intern с БД и API соответственно. У каждого контейнера есть свой Dockerfile.
API будет запущено на `localhost:8080`, а Tarantool на `tarantool:3301`. Порт и адресс, на котором запускается Tarantool можно поменять в `./config/config.yml`.

## 4. API Эндпоинты

При попытке отправки неккоректного ответа и неккоректного чтения запроса json формата во всех эндпоинтах предусмотрен хендлинг ошибок.

### 1. Аутентификация (POST /api/login)
Этот эндпоинт выполняет аутентификацию пользователя и возвращает JWT-токен.<br>
Можно повторно авторизоваться (в таком случае, токен сессии в БД для данного юзера перезапишется, обновится).
Пример запроса: </br>
```json
POST /api/login HTTP/1.1
Host: localhost:8080
Content-Type: application/json
{
  "username": "admin",
  "password": "presale"
}
```
Ответ, статус `200 Ok`: 
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiYWRtaW4iLCJleHAiOjE3MjQxMTExNDV9.0kn4u7X-JhO-eHZf8IOc_zWKONL42-tvFMbKkD1fibo"
}
```

При отсылке некорректных данных получим ошибку и статус `401 Unauthorized`:
```json
{
    "status": "Invalid username or password"
}
```
При ошибке создании (перезаписи) сессии в БД, получим ошибку и статус `401 Unauthorized`:
```json
{
    "status": "Failed to create (rewrite) session"
}
```
### 1. Запись (перезапись) данных (POST /api/write)
Этот эндпоинт позволяет авторизованным пользователям записывать пары ключ-значение в базу данных Tarantool, причем ключ только типа `string`, значение только `scalar (bool, int, float, string)`.
Запись происходит асинхронно, с помощью батчей (реализованы вручную в `./internal/logic/data.go`).
Пример запроса:
```json
POST /api/login HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiYWRtaW4iLCJleHAiOjE3MjQxMTExNDV9.0kn4u7X-JhO-eHZf8IOc_zWKONL42-tvFMbKkD1fibo
Content-Type: application/json
{
  "data": {
    "key1": "This is a string",
    "key2": 321,
    "key3": true,
    "key4": 123.4567890123456789
  }
}
```
Ответ, статус `200 Ok`:
```json
{
    "status": "Success"
}
```
При попытке отправки неккоректного JWT токена получим ошибку, статус `401 Unauthorized`:
```json
{
    "status": "Not authorized: signature is invalid"
}
```
При попытке отправки неккоректного Authorization заголовка получим ошибку, статус `401 Unauthorized`:
```json
{
    "status": "Not authorized, wrong header format"
}
```
При попытке отправки просроченого JWT токена получим ошнибку, статус `401 Unauthorized`  (просроченная сессия удаляется из БД):
```json
{
    "status": "Not authorized: session expired"
}
```
При попытке отправки несуществующей (или же уже удаленной) сессии получим ошибку, статус `401 Unauthorized`:
```json
{
    "status": "Not authorized: session not found"
}
```
При попытке отправки плохих данных (не скалярного типа) получим ошибку, статус `400 Bad Request`.
При этом данные в БД никакие на запишутся, не перезапишутся за счет валидации данных перед записью.
Пример запроса:
```json
{
  "data": {
    "key1": "This is a string",
    "key2": 321,
    "key3": true,
    "key4": 123.4567890123456789,
    "key5": [123]
  }
}
```
Ответ:
```json
{
    "status": "Bad Request: invalid value type for key key5 (only scalar values are allowed)"
}
```
При ошибке записи в БД (чего добиться скорее всего не получится), выдастся ошибка, статус `400, Bad Request`:
```json
{
    "status": "Failed to write batch. Errors: (key: 'key1', error: 'error1')"
}
```

  




