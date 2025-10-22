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
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
```


