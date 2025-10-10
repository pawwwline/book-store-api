# BookStore API


> REST-сервис на Go для управления книгами с PostgreSQL и кэшем для ускоренного доступа. Позволяет добавлять, получать, обновлять и удалять книги. Документация доступна через Swagger.


## Требования
Убедитесь, что у вас установлен docker и docker compose, go v1.25
```bash
docker --version
docker compose version
go version
```

## Установка и запуск

1. Клонировать репозиторий
```bash
git clone https://github.com/pawwwline/book-store-api
cd book-store-api
```

2. Создать .env на основе .env.example (при необходимости поставить нужные значения)
```bash
cp env.example .env
```
> **Note:** Доступные APP_ENV `local` `test` `dev` `prod`

3. Запуск сервиса

```bash
make run
```
>**Note:** Автоматически запустит миграции

### Миграции

Для применения миграций

```bash
make migrate-up
```

Для отката миграций

```bash
make migrate-up
```



### Документация

> **Note:** Документация API доступна на `api/v1/swagger/index.html` после запуска сервиса.


