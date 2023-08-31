## Какие возможности имеет программа ##
1. Создание и удаление сегментов;
2. Удаление старых и добавление новых сегментов пользователю;
3. Получение актуальных сегментов пользователя;
4. Получение истории изменений сегментов пользователя в формате CSV;
5. Добавление пользователей в сегмент на время.

## Вопросы, возникшие во время разработки ##
### Хранение сегментов и истории ###
В связи с необходимостью поддерживать целостность истории операций, при удалении существующего сегмента лишь устанавливается метка о его неактивности.

### Возможности подробного анализа операций над сегментами, в которые был добавлений пользователь ###
Были добавления новые операции (`обновлено значение`, `удален в связи с удалением сегмента` и т.д.), позволяющие более четко контролировать изменения состояния пользователей.


## Инструкции по запуску ##

1. Запуск через docker-compose. В базу данных автоматически загрузится необходимая схема и таблицы. Сервер будет доступен по порту `8080`, а база данных - `8001`.
2. Запуск бинарного файла. Настройки по умолчанию можно найти в файле `/config/config-example.yml`. Для данного случая необходимо поставить `Postgres версии 15` и выполнить скрипт инициализации `init.sql`. 

После запуска программа выведен сообщение о начале работы. С этого момента можно начинать отправлять запросы через HTTP API.

Скрипт инициализации базы данных не содержит пользователей или сегментов. Сегменты добавляются вручную, пользователи - автоматически при вызове добавления в сегмент.

Приоритет конфигурации: `config.yml рядом с бинарником` > `переменные окружения` > `настройки по умолчанию`.

## Примеры запросов ##

`Swagger`-документацию по API можно найти в папке `/docs`, а также по адресу `/swagger` включенного сервиса.

### Создание сегмента ###
#### Запрос ####
`PUT /segment`

```
curl --location --request PUT 'localhost:8080/segment' \
--header 'Content-Type: application/json' \
--data '{
    "slug" : "TEST_SEG"
}'
```

### Удаление сегмента ###
#### Запрос ####
`DELETE /segment`

```
curl --location --request DELETE 'localhost:8080/segment' \
--header 'Content-Type: application/json' \
--data '{
    "slug" : "TEST_SEG_3"
}'
```

### Добавление сегмента 10% пользователям до 10 сентября 2023 ###
#### Запрос ####
`PUT /segment`

```
curl --location --request PUT 'localhost:8080/segment' \
--header 'Content-Type: application/json' \
--data '{
    "slug" : "TEST_SEG_3",
    "untilDate" : "2023-09-10T13:29:07+03:00",
    "percent" : "10"
}'
```

### Простое добавление пользователя в сегмент ###
#### Запрос ####
`PUT /assignment`

```
curl --location --request PUT 'localhost:8080/assignment' \
--header 'Content-Type: application/json' \
--data '{
    "userID" : "1035",
    "segmentAdd" : [
        {
            "slug" : "TEST_SEG"
        }
    ],
    "segmentRemove" : []
}'
```

### Обновление TTL сегмента пользователя ###
#### Запрос ####
`PUT /assignment`

```
curl --location --request PUT 'localhost:8080/assignment' \
--header 'Content-Type: application/json' \
--data '{
    "userID" : "1026",
    "segmentAdd" : [
        {
            "slug" : "TEST_SEG",
            "untilDate" : "2023-09-05T06:51:13+03:00"
        }
    ],
    "segmentRemove" : []
}'
```

### Удаление из сегмента и добавление в новый ###
#### Запрос ####
`PUT /assignment`

```
curl --location --request PUT 'localhost:8080/assignment' \
--header 'Content-Type: application/json' \
--data '{
    "userID" : "1026",
    "segmentAdd" : [
        {
            "slug" : "TEST_SEG_2"
        }
    ],
    "segmentRemove" : [
        "TEST_SEG"
    ]
}'
```

### Получение списка текущих сегментов пользователя ###
#### Запрос ####
`GET /user`

```
curl --location --request GET 'localhost:8080/user' \
--header 'Content-Type: application/json' \
--data '{
    "userID" : "1026"
}'
```
#### Ответ ####
```
    HTTP/1.1 200 OK
    Date: Thu, 31 Aug 2023 14:21:05 GMT
    Status: 200 OK
    Content-Type: application/json
    Content-Length: 140


{
    "userID": 1035,
    "segmentAdd": [
        {
            "slug": "TEST_SEG_2",
            "untilDate": ""
        },
        {
            "slug": "TEST_SEG_3",
            "untilDate": ""
        },
        {
            "slug": "TEST_SEG",
            "untilDate": ""
        }
    ]
}
```
### Получение истории изменений сегментов пользователя в формате CSV ###
#### Запрос ####
`GET /user/history`

```
curl --location --request GET 'localhost:8080/user/history' \
--header 'Content-Type: application/json' \
--data '{
    "userID" : "1026",
    "year" : "2023",
    "month" : "8"
}'
```
#### Ответ ####
`200 OK`
```
HTTP/1.1 200 OK
Date: Thu, 31 Aug 2023 14:18:09 GMT
Status: 200 OK
Content-Type: text/csv
Content-Length: 619

1035;TEST_SEG;Добавление;2023-08-31T17:16:05+03:00
1035;TEST_SEG_2;Добавление;2023-08-31T17:16:05+03:00
1035;TEST_SEG_3;Добавление;2023-08-31T17:16:05+03:00
1035;TEST_SEG;Обновление значения;2023-08-31T17:16:20+03:00
1035;TEST_SEG;Добавление;2023-08-31T17:17:44+03:00
1035;TEST_SEG;Удаление;2023-08-31T17:17:44+03:00
1035;TEST_SEG;Добавление;2023-08-31T17:17:57+03:00
1035;TEST_SEG;Удаление;2023-08-31T17:17:57+03:00
1035;TEST_SEG;Добавление;2023-08-31T17:18:00+03:00
1035;TEST_SEG;Удаление;2023-08-31T17:18:00+03:00
```
## Внешние зависимости ##
1. go-chi/chi - маршрутизация, middleware;
2. rs/zerolog - логгирование;
3. ilyakaznacheev/cleanenv - загрузка конфигурации;
4. go-playground/validator/v10 - валидация данных;
5. pgx/v5 - драйвер Postgres;
6. scany/v5 - чтение данных из БД в структуры;
7. swaggo/swag - генерация документации API.