# Веб-сервис Арифметических Вычислений
## Описание
Этот проект реализует веб-сервис для вычисления арифметических выражений. Пользователи могут отправлять математические выражения через HTTP-запросы, а сервис возвращает результат вычисления. Сервис поддерживает базовые арифметические операции: сложение, вычитание, умножение, деление, а также использование скобок для определения приоритетов операций.
## Структура Проекта
```
.
├── cmd
│   └── main.go              # Точка входа приложения
├── internal
│   ├── calc
│   │   ├── calc.go          # Логика вычисления выражений
│   │   └── calc_test.go     # Тесты для калькулятора
│   └── server
│       ├── api.go           # Реализация эндпоинта api/v1/calculate
│       ├── server.go        # HTTP-сервер
│       └── server_test.go   # Тесты для сервера
├── go.mod
```
## Установка
1. **Клонируйте Репозиторий:**
```bash
   git clone https://github.com/Yorshik/final_task_sprint_1.git
```
2. **Перейдите в Каталог Проекта:**
```bash
   cd final_task_sprint_1
```
3. **Установите Необходимые Зависимости:**
```bash
   go mod download
```
## Запуск Проекта
Для запуска веб-сервиса выполните следующую команду из корневого каталога проекта:
```bash
go run ./cmd/main.go/
```
Сервис будет доступен по адресу [localhost:9000/api/v1/calculate](http://localhost:9000/api/v1/calculate).
## Использование
### Шаблон Запроса
```bash
curl -i -X POST http://localhost:9000/api/v1/calculate/ \
     -H "Content-Type: application/json" \
     -d '{"expression":"EXPRESSION"}'
```
### Важно
- **Только POST-запросы** позволяют получить результат вычисления выражения.
- **Внимание на Windows:** При использовании `curl` на Windows необходимо экранировать двойные кавычки в части `-d '{"expression":"EXPRESSION"}'`, заменив одинарные кавычки на двойные.
### 1. Успешные Запросы
#### Пример 1: Простое Выражение
**Команда `curl`:**
```bash
curl -i -X POST http://localhost:9000/api/v1/calculate/ \
     -H "Content-Type: application/json" \
     -d '{"expression":"2+2*2"}'
```
**Ожидаемый Ответ:**
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 01 Jan 1970 00:00:00 GMT
Content-Length: 13
{
  "result": 6
}
```
#### Пример 2: Выражение со Скобками
**Команда `curl`:**
```bash
curl -i -X POST http://localhost:9000/api/v1/calculate/ \
     -H "Content-Type: application/json" \
     -d '{"expression":"(3 + 4) * 5"}'
```
**Ожидаемый Ответ:**
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 01 Jan 1970 00:00:00 GMT
Content-Length: 14
{
  "result": 35
}
```
#### Пример 3: Выражение с Пробелами
**Команда `curl`:**
```bash
curl -i -X POST http://localhost:9000/api/v1/calculate/ \
     -H "Content-Type: application/json" \
     -d '{"expression":" 10 / (2 + 3) "}'
```
**Ожидаемый Ответ:**

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 01 Jan 1970 00:00:00 GMT
Content-Length: 15
{
  "result": 2
}
```
### 2. Ошибки 422 - Некорректные Входные Данные
#### Случай 1: Содержится Буква
**Команда `curl`:**
```bash
curl -i -X POST http://localhost:9000/api/v1/calculate/ \
     -H "Content-Type: application/json" \
     -d '{"expression":"2+2*2a"}'
```
**Ожидаемый Ответ:**
```
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json
Date: Thu, 01 Jan 1970 00:00:00 GMT
Content-Length: 34
{
  "error": "Expression is not valid"
}
```
#### Случай 2: Неправильное количество скобок и/или неправильная постановка скобок
**Команда `curl`:**
```bash
curl -i -X POST http://localhost:9000/api/v1/calculate/ \
     -H "Content-Type: application/json" \
     -d '{"expression":"(3 + 4 * 2"}'
```
**Ожидаемый Ответ:**
```
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json
Date: Thu, 01 Jan 1970 00:00:00 GMT
Content-Length: 34
{
  "error": "Expression is not valid"
}
```
#### Случай 3: Деление на Ноль
**Команда `curl`:**
```bash
curl -i -X POST http://localhost:9000/api/v1/calculate/ \
     -H "Content-Type: application/json" \
     -d '{"expression":"10/0"}'
```
**Ожидаемый Ответ:**
```
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json
Date: Thu, 01 Jan 1970 00:00:00 GMT
Content-Length: 34
{
  "error": "Expression is not valid"
}
```
### Предупреждения
- **Внутренние Ошибки (500):** 

В некоторых случаях может возникнуть внутренняя ошибка сервера, например, из-за непредвиденных сбоев. При таких ошибках вы получите HTTP-ответ:
```
HTTP/1.1 500 Internal Server Error
  Content-Type: application/json
  Date: Thu, 01 Jan 1970 00:00:00 GMT
  Content-Length: 34
{
    "error": "Internal server error"
  }
```
- **Использование `curl` на Windows:**

При выполнении команд `curl` на Windows необходимо экранировать двойные кавычки и использовать двойные кавычки для всего JSON-объекта. Например:
```commandline
  curl -i -X POST http://localhost:9000/api/v1/calculate/ ^
       -H "Content-Type: application/json" ^
       -d "{\"expression\":\"2+2*2\"}"
```
Обратите внимание на использование `^` для переноса строк в командной строке Windows и экранирование внутренних двойных кавычек с помощью \.
## Тестирование
Для запуска тестов используйте команду:
```bash
go test ./...
```
Это запустит все тесты в проекте и отобразит результаты.
