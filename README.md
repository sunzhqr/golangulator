
# Golangulator — Telegram-бот калькулятор на Go

![Go](https://img.shields.io/badge/Go-1.22-blue) ![Telegram](https://img.shields.io/badge/Telegram-Bot-green) ![GORM](https://img.shields.io/badge/GORM-PostgreSQL-lightgrey) ![Zap](https://img.shields.io/badge/Logging-Zap-orange)

---

## Описание

Golangulator — это Telegram-бот, который умеет вычислять арифметические выражения, сохранять историю вычислений пользователей в PostgreSQL и возвращать её по запросу.  
Поддерживается:

- вычисление выражений с операциями +, -, *, /, **, ^, % скобками
- сохранение истории вычислений пользователя
- очистка истории
- inline режим — вычисление выражения прямо в чате через inline-запросы
- можно включить и отключить сохранения истории inline-запроса
- логирование через Uber Zap
- база данных PostgreSQL для хранения истории через GORM ORM

---

## Функциональность

- `/start` — приветственное сообщение с инструкциями  
- `/help` — помощь по использованию бота  
- `/history` — показать последние вычисления пользователя  
- `/clear_history` — очистить историю вычислений  
- Вычисление арифметических выражений из сообщений  
- Inline-режим для быстрого вычисления без отправки сообщения  

---

## Технологии

- Go 1.22+  
- Telegram Bot API v5 (github.com/go-telegram-bot-api/telegram-bot-api)  
- GORM ORM для PostgreSQL  
- PostgreSQL (для хранения истории вычислений)  
- Uber Zap для логирования  
- dotenv для загрузки конфигурации из `.env`

---

## Установка и запуск

### Шаг 1: Клонирование репозитория

```bash
git clone https://github.com/sunzhqr/golangulator.git
cd golangulator
```

### Шаг 2: Настройка `.env` файла

Создайте файл `.env` в корне проекта:

```
BOT_TOKEN=ваш_токен_бота_в_telegram
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

### Шаг 3: Запуск PostgreSQL

Убедитесь, что у вас запущена база PostgreSQL, и указанные в DATABASE_URL данные верны.

### Шаг 4: Установка зависимостей и запуск

```bash
go mod download
go run cmd/bot/main.go
```

---

## Структура проекта

```
├── cmd
│   └── bot
│       └── main.go            # Точка входа
├── internal
│   ├── domain                # Интерфейсы и доменная логика
│   ├── infrastructure
│   │   ├── storage           # Репозиторий для работы с БД
│   │   └── telegram          # Логика Telegram-бота
│   └── usecase               # Бизнес-логика: калькулятор, история
├── go.mod
├── go.sum
└── .env
```

---

## Пример использования

Отправьте боту сообщение с выражением, например:

```
2 + 2 * (3 - 1)
```

Бот ответит:

```
Результат: 6
```

Чтобы посмотреть историю, отправьте команду:

```
/history
```

---

## Логирование

Бот использует Uber Zap для структурированного логирования. Логи по умолчанию выводятся в консоль.

---

@sunzhqr all rights reserved)
