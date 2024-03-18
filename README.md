# warehouse

- Описание API находится по адресу `http://localhost:8080/swagger/`
- Коллекция `Postman` лежит в директории `postman`

## Запуск сервиса

Необходимо иметь локально установленные утилиты

- make
- docker

Находясь в корне проекта:

```shell
make up
```

Проверка успешного старта:

```shell
curl localhost:8080/ping
```

Ответ должен иметь следующий вид:

```shell
{"status":"ok","message":""}
```
