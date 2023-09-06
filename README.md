# l0wb

# Развертывание и настройка сервиса заказов

## Шаг 1: Установка и настройка PostgreSQL

- <input type="checkbox" checked="checked" disabled="disabled"> Установить PostgreSQL локально.
- <input type="checkbox" checked="checked" disabled="disabled"> Создать новую базу данных для заказов.
- <input type="checkbox" checked="checked" disabled="disabled"> Создать пользователя с доступом к базе данных.
- Настроить доступ к PostgreSQL для приложения.

## Шаг 2: Создание таблиц для хранения данных заказов

- <input type="checkbox" checked="checked" disabled="disabled"> Определить структуру таблицы для заказов и связанных данных.
- <input type="checkbox" checked="checked" disabled="disabled"> Создать таблицы в базе данных с помощью SQL-запросов.

## Шаг 3: Разработка сервиса заказов

- <input type="checkbox" checked="checked" disabled="disabled"> Создать Go-приложение для сервиса заказов.
- <input type="checkbox" checked="checked" disabled="disabled"> Реализовать подключение и подписку на канал в NATS Streaming.
- <input type="checkbox" checked="checked" disabled="disabled"> Записать полученные данные из канала в базу данных PostgreSQL.

## Шаг 4: Реализация кэширования данных

- <input type="checkbox" checked="checked" disabled="disabled"> Создать кэш для данных заказов в приложении.
- <input type="checkbox" checked="checked" disabled="disabled"> Добавить механизм очистки кэша от истекших данных.
- <input type="checkbox" checked="checked" disabled="disabled"> Восстановить кэш из базы данных при падении сервиса.

## Шаг 5: Разработка HTTP-сервера

- <input type="checkbox" checked="checked" disabled="disabled"> Разработать HTTP-сервер для доступа к данным заказов.
- <input type="checkbox" checked="checked" disabled="disabled"> Создать API-методы для отображения и поиска заказов.

## Шаг 6: Запуск сервиса

- <input type="checkbox" checked="checked" disabled="disabled"> Запустить сервис заказов на локальной машине.
