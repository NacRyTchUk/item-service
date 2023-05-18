# item-service

__Golang Микросервис реализующий CRUD методы работы с ITEMS__
* REST [API](https://nrtu.stoplight.io/docs/item-service) через http gateway, с возможнотью добавления gRPC как транспорт
* Хранение данных в Postgres, кеширование в Redis
* Логи операций хранятся в Clickhouse, передающиеся через Nats

Проект является реализацией [тестового задания]()

# setup

```bash
$ docker-compose up --build # Для запуска всего стека

$ make run                  # Для запуска только микросервиса
$ make run-env              # Для запуска Nats+Clickhouse
```