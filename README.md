# MIFI Bank - REST API для банковского сервиса

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-blue)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-green)](https://www.postgresql.org/)

Проект для дисциплины "Язык программирования GO". Магистратура МИФИ. Специальность "Программная инженерия". 2-й семестр.

## Особенности

- 🔐 Регистрация и аутентификация с JWT
- 💳 Генерация карт по алгоритму Луна
- ↔️ Переводы между счетами с проверкой баланса
- 🏦 Кредитные операции с аннуитетными платежами
- 📊 Финансовая аналитика операций
- 📧 Email-уведомления через SMTP
- 🔒 Шифрование данных (PGP, HMAC-SHA256)

## Технологии

- **Язык:** Go 1.23+
- **База данных:** PostgreSQL 17 (с расширением pgcrypto)
- **Библиотеки:**
    - Маршрутизация: `gorilla/mux`
    - Аутентификация: `golang-jwt/jwt/v5`
    - Логирование: `logrus`
    - Шифрование: `bcrypt`, `crypto/hmac`
    - SMTP: `gomail.v2`

## Примеры запросов

### Регистрация:
```bash
curl -X POST http://localhost:8080/register \
-H "Content-Type: application/json" \
-d '{"email":"user@example.com", "username":"john_doe", "password":"securePass123"}'
```
### Аутентификация:
```bash
curl -X POST http://localhost:8080/login \
-H "Content-Type: application/json" \
-d '{"email":"user@example.com", "password":"securePass123"}'
```
### Создание счета (требует JWT):
```bash
curl -X POST http://localhost:8080/accounts \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-H "Content-Type: application/json" \
-d '{"currency":"RUB"}'
```

## Эндпоинты API
### Публичные эндпоинты

| Метод | Путь       | Описание                  |
|-------|------------|---------------------------|
| POST  | /register  | Регистрация пользователя  |
| POST  | /login     | Аутентификация           |

### Защищенные эндпоинты (JWT)

| Метод | Путь                              | Описание                          |
|-------|-----------------------------------|-----------------------------------|
| POST  | /accounts                         | Создать новый счет               |
| POST  | /cards                            | Выпустить виртуальную карту       |
| POST  | /transfer                         | Перевод между счетами            |
| GET   | /analytics                        | Получить финансовую аналитику    |
| POST  | /credits/create                   | Оформить кредит                  |
| GET   | /credits/{creditId}/schedule      | График платежей по кредиту       |

## Безопасность

    🔑 Хеширование паролей с bcrypt

    🛡️ JWT-токены с 24-часовым сроком действия

    🔒 Шифрование данных карт (PGP + HMAC)

    🛡️ Защита от SQL-инъекций (параметризованные запросы)

    📨 HTTPS рекомендуется для продакшена

## Лицензия

- MIT License.
