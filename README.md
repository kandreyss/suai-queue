# 🎓 SUAI Queue Bot

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Telegram](https://img.shields.io/badge/Telebot-v3-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)
![Gorm](https://img.shields.io/badge/ORM-GORM-red?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**SUAI Queue Bot** — это Telegram-бот для студентов ГУАП, предназначенный для управления очередью на сдачу лабораторных работ и зачетов.

Бот позволяет студентам записываться в электронную очередь, отслеживать свое место в реальном времени и получать уведомления. Проект построен на языке Go с использованием чистой архитектуры (Clean Architecture).

## ✨ Возможности

- 📝 **Регистрация**: Привязка Telegram-аккаунта к имени и учебной группе.
- 🔄 **Живая очередь**:
  - Запись в очередь одной кнопкой.
  - Просмотр текущего списка студентов.
  - Выход из очереди после сдачи.
- ⚙️ **Управление профилем**: Возможность сменить имя или номер группы.
- 🧹 **Автоматическая очистка**:
  - Бот отслеживает время нахождения в очереди.
  - Удаляет "мертвые души" (студентов, которые забыли выйти) спустя заданное время (по умолчанию 25 минут).
  - Отправляет уведомление удаленному пользователю.
- 🗄 **База данных**: Хранение пользователей и их данных в SQLite.

## 🛠 Технологический стек

- **Язык**: Go (Golang)
- **Архитектура**: Standard Go Project Layout + Clean Architecture
- **Telegram Bot API**: [gopkg.in/telebot.v3](https://gopkg.in/telebot.v3)
- **Database**: SQLite (через [GORM](https://gorm.io/))
- **Конфигурация**: [godotenv](https://github.com/joho/godotenv)

## 📂 Структура проекта

Проект следует стандарту **Standard Go Project Layout**:

```
suai-queue/
├── cmd/
│   └── bot/
│       └── main.go              # Точка входа в приложение
├── internal/
│   ├── app/                     # Инициализация и сборка зависимостей (DI)
│   ├── config/                  # Загрузка конфигурации
│   ├── domain/                  # Чистые бизнес-сущности (Student)
│   ├── repository/              # Слой данных (работа с SQLite через GORM)
│   ├── service/                 # Бизнес-логика
│   │   └── queue/               # Логика очереди и фоновая очистка
│   ├── session/                 # Управление сессиями (FSM)
│   └── transport/
│       └── telegram/            # Хендлеры и маршрутизация Telegram
├── storage/                     # Файл базы данных (SQLite)
├── .env                         # Переменные окружения
└── Makefile                     # Команды для сборки и запуска
```

## 🚀 Установка и запуск

### Предварительные требования
- Go 1.22 или выше
- Утилита `make` (опционально)
- Токен Telegram бота (от [@BotFather](https://t.me/BotFather))

### 1. Клонирование репозитория
```bash
git clone https://github.com/your-username/suai-queue.git
cd suai-queue
```

### 2. Настройка окружения
Создайте файл `.env` в корне проекта:
```bash
TOKEN=your_telegram_bot_token
DB_PATH=storage/suai_queue.db
```

### 3. Запуск
Используя Makefile:
```bash
make run
```
Или напрямую через Go:
```bash
go run cmd/bot/main.go
```

### 4. Сборка бинарного файла
```bash
make build
# Запуск собранного файла
./bin/suai-queue
```

## 🎮 Использование

1. Найдите бота в Telegram и нажмите **Start**.
2. Зарегистрируйтесь командой `/register` (введите имя и номер группы).
3. Используйте меню:
   - **➕ Встать в очередь**: Занять место.
   - **📋 Посмотреть очередь**: Узнать, кто сейчас сдает.
   - **➖ Выйти из очереди**: Освободить место.
   - **ℹ️ Инфо**: Посмотреть свои данные.

## 🤝 Вклад в развитие (Contributing)

Мы приветствуем Pull Request'ы! Если вы хотите улучшить проект:
1. Форкните репозиторий.
2. Создайте ветку (`git checkout -b feature/amazing-feature`).
3. Внесите изменения и закоммитьте их.
4. Откройте Pull Request.

## 📄 Лицензия

Distributed under the MIT License. See `LICENSE` for more information.
