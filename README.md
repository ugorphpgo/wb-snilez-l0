# #Задание l0 - горутиновый golang

Микросервис, который получает данные о заказах из Kafka, сохраняет их в Postgres, кэширует в памяти и предоставляет HTTP API и веб-интерфейс для просмотра информации по id заказа (`order_uid`).

---

## 
- Подписка на Kafka (топик `orders`)  
- Валидация и парсинг JSON сообщений  
- Сохранение заказов в PostgreSQL (UPSERT по `order_uid`)  
- Кэширование заказов в памяти для быстрого доступа  
- Восстановление кэша из БД при запуске  
- HTTP API:
  - `GET /order/{order_uid}` — возвращает заказ в формате JSON  
- Веб-страница:
  - `GET /` — форма для поиска заказа по `order_uid`  

---

## С помощью каких технологий реализован проект
- **Go 1.25**  
- **Kafka**  
- **PostgreSQL** 
- **Docker** — контейнеризация Kafka и Postgres и запуск программы локально c помощью docker-compose
- **Библиотеки, которые использовал в проекте**:  
  `segmentio/kafka-go`, `pgx/v5`, `zap`, `viper`, `golang-migrate`

---

## Quick start

1. Запустить окружение (БД + Kafka + сервис):
   ```bash
   docker compose up -d --build
   ```

2. В отдельном терминале запускаем продюсер(генератор заказов):
   ```bash
   go run ./cmd/producer
   ```

3. Открыть веб-интерфейс:
   [http://localhost:8081](http://localhost:8081)

---

##  Структура проекта
```
cmd/
  producer/      — отправка тестовых сообщений в Kafka
  wbservice/     — основной HTTP-сервис
configs/
  config.yaml    — конфигурация сервиса
internal/
  app/           — инициализация приложения
  cache/         — реализация кэша в памяти
  config/        — загрузка конфигурации
  http/          — обработчики и сервер
  kafka/         — получение сообщений из Kafka
  log/           — логирование (zap)
  model/         — модель данных заказа
  repo/          — работа с PostgreSQL и миграции
  service/       — бизнес-логика
migrations/      — SQL миграции
web/
  index.html     — веб-интерфейс для поиска заказа
```

---

##  Пример запроса
```bash
GET http://localhost:8081/order/b563feb7b2b84b6test
```
Ответ:
```json
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "delivery": { "name": "Test Testov", "city": "Kiryat Mozkin" },
  "payment": { "amount": 1817, "currency": "USD" },
  "items": [{ "name": "Mascaras", "price": 453 }]
}
```


