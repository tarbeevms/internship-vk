# Документация к API

## 1. Обзор

Этот API предназначен для работы с базой данных Tarantool, обработки пользовательской аутентификации, управления сессиями и хранения данных в формате ключ-значение.<br>
Для аутентификации используется JWT-токен. У токена есть время действия (выдается на 3 минуты), по истечению сессия удаляется и необходимо заново авторизоваться. Регистрация не предусмотрена, существует только один пользователь по умолчанию (`username: admin, password: presale`).<br>
Сервис предоставляет API для входа, записи и чтения данных, соответственно по `/api/login`, `/api/write`, `/api/read`.<br>

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

### 1. Аутентификация (POST /api/login)
Этот эндпоинт выполняет аутентификацию пользователя и возвращает JWT-токен.<br>
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
Ответ (статус `200 Ok`): 
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
При попытке авторизоваться за того же пользователя повторно (либо же при ошибке создании сесси в БД, но сломать у вас все равно не получится с:), получим ошибку и статус `401 Unauthorized`:
```json
{
    "status": "Failed to create session or session already exists"
}
```
### 1. Аутентификация (POST /api/login)





