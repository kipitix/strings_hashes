# StringsHashes

 GO: Final project

---

## Задание

### Требуется реализовать следующий набор сервисов:

- [✔️] Service 1: Stateless сервис, который принимает на вход набор произвольных строк, считает от них хэши SHA-3 и возвращает вызывающей стороне.
- [✔️] Service 2: Stateful сервис, который соответствует следующий спецификации (Swagger расположен в конце данного текста).

---

### Требования к реализации

К Git репозиторию:

- [✔️] Репозиторий должен храниться в приватном репозитории в gitlab.rebrainme.com в группе gitlab.rebrainme.com/golang_users_repos/<your_gitlab_id>/

К БД:

- [✔️] Выбор БД для хранения хэшей за вами

К Сервисам:

- [✔️] Для межсервисного взаимодействия используйте gRPC
- [✔️] Выбор способа хранения конфигураций (данные для соединения с бд, адрес соседнего сервиса и т.д.) за вами
- [✔️] Сервисы должны логгировать свое поведение
- [✔️] Service 2 должен реализовывать спецификацию,описанную в swagger’е
- [✔️] Service 1 подсчет хэшей должен быть организован параллельно
- [✔️] Service 1 должен быть покрыт тестами
- [✔️] Оба сервиса должны находиться в контейнерах и подниматься в связке с помощью docker-compose
- [✔️] Оба сервиса должны корректно обрабатывать сигналы о завершении работы (graceful shutdown)

К логам:

- [✔️] Сервисы должны писать логи в формате GELF
- [✔️] Логирование должно быть сквозным
- [✔️] К логам должен быть добавлено поле со Stack Trace’ом. Для получения трейса можно использовать:

```go
    import "github.com/pkg/errors"

    errors.WithStack(err))
```

- [✔️] Продумайте использование уровней логирования (вам надо будет аргументировать во время презентации ваше решение)

---

### Порядок действий во время презентации проекта:

- [ ] Пришлите в ответ на это сообщение ссылку на ваш репозиторий, а так же предоставьте ее куратору.
- [ ] Продемонстрируйте компиляцию и запуск связки сервисов
- [ ] Продемонстрируйте работу сервисов и написанные тесты
- [ ] Ответьте на заданные вопросы куратора по проекту, объясните выбор некоторых решений
- [ ] Ответьте на заданные вопросы куратора по всему практикуму

---

### Swagger для Service 2

```yaml
swagger: "2.0"
info:
  version: "1.0.0"
  title: "Итоговое задание. Хэши."
  description: "Данный сервис должен, взаимодействуя с сервисом считающим хэши (по выбранному вами протоколу), получать из входящих строк их хэши, сохранять их в свою БД (выбор так же за вами) с присвоением id, по которым далее можно будет запрашивать хэши."
schemes:
- "http"
produces:
  - application/json
paths:
  /send:
    post:
      summary: "Получает на вход список строк, хэши от которых нужно посчитать и сохранить"
      parameters:
        - in: body
          name: params
          description: "Strings for hash"
          schema:
            $ref: '#/definitions/ArrayOfStrings'
      responses:
        "200":
          description: "Success"
          schema:
            $ref: '#/definitions/ArrayOfHash'
        "400":
          description: "Bad request"
        "500":
          description: "Internal Server Error"
  /check:
    get:
      summary: "Получает по id хэш из хранилища (если есть)"
      parameters:
        - in: query
          name: ids
          description: "Get hash by this id"
          required: true
          type: array
          items:
            type: string
      responses:
        "200":
          description: "Success"
          schema:
            $ref: '#/definitions/ArrayOfHash'
        "204":
          description: "No Content"
        "400":
          description: "Bad request"
        "500":
          description: "Internal Server Error"
definitions:
  ArrayOfStrings:
    type: array
    items:
      type: string
  ArrayOfHash:
    type: array
    items:
      $ref: '#/definitions/Hash'
  Hash:
    type: object
    properties:
      id:
        type: integer
        example: 38
      hash:
        type: string
        example: a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a
    required:
      - id
      - hash
```

---

## Пояснения к реализации

### Требуется реализовать следующий набор сервисов:

#### Service 1: Stateless сервис, который принимает на вход набор произвольных строк, считает от них хэши SHA-3 и возвращает вызывающей стороне.

Service 1 - stateless сервис, который вычисляет хеши по SHA-3.

В репозитории это go модуль [hashcalc](./hashcalc/).

Сервис необходим только для вычислений, поэтому у него простая архитектура - gRPC сервер и пакет для вычислений.

Описание gRPC находится в файле [hashcalc.proto](./hashcalc/api/hashcalc.proto).

Для обеспечения параллельности вычислений используются каналы.

Данные вычитываются из gRPC stream, подаются во входной канал вычислителя контрольных сумм, отдаются из его выходного канала и результат отправляется по выходному gRPC stream.

---

#### Service 2: Stateful сервис, который соответствует следующий спецификации (Swagger расположен в конце данного текста).

Service 2 - stateful сервис, который обрабатывает запросы от клиента.

В репозитории это go модуль [hashkeeper](./hashkeeper/).

---

##### Структура модуля

В каталоге [internal](./hashkeeper/internal/) основные компоненты сервиса (пакеты разделены в соответствии с подходом к чистой архитектуре):
- [domain](./hashkeeper/internal/domain/) - сущности для обеспечения логики работы (хеш, пайплайн вычисления хешей и пайплайн запроса хешей)
  - [datahash](./hashkeeper/internal/domain/datahash/) - основная сущность хеша `Hash`, интерфейс создателя хешей `HashMaker`, интерфейс репозитория хешей `HashRepository`
  - [calculate](./hashkeeper/internal/domain/calculate/) - пайплайн для вычисления и сохранения хешей, использует `HashMaker` и `HashRepository`
  - [find](./hashkeeper/internal/domain/find/) - пайплайн для запроса хешей из базы, использует `HashRepository`
- [application](./hashkeeper/internal/application/) - приложение - сущность для передачи собранных во едино компонентов
- [infrastructure](./hashkeeper/internal/infrastructure/) - инфраструктура - реализация репозитория для конкретной базы данных
- [interfaces](./hashkeeper/internal/interfaces/) - интерфейсы - реализация REST API сервера и gRPC клиента

---

##### Логика работы

У сервиса 2 ручки: `/send` и `/check`

Последовательность действий для ручки `/send`:
- Получаем запрос
- Вычисляем хеши
- Записываем хеши в базу данных
- Запрашиваем хеши из базы по содержимому хешей, получая информацию об идентификаторах
- Отправляем ответ клиенту

*Важный момент! Когда записываем хеши в базу, то идентификаторы образуются за счёт первичного ключа базы. Так же в базе есть ограничение на колонку с хешем и при записи дубликатов в базе не будет. После записи производится запрос хешей из базы по содержимому хешей. Таким образом можно избежать дубликатов одинаковых хешей и выдавать клиенту идентификаторы уже имеющихся записей.*

Последовательность действий для ручки `/check`:
- Получаем запрос
- Запрашиваем хеши из базу данных по ключам
- Отправляем ответ клиенту

---

### Соответствие требований и реализации

### К Git репозиторию:

#### Репозиторий должен храниться в приватном репозитории в gitlab.rebrainme.com в группе gitlab.rebrainme.com/golang_users_repos/<your_gitlab_id>/

[Ссылка на репозиторий](https://gitlab.rebrainme.com/golang_users_repos/3439/stringshashes)

---

### К БД:

#### Выбор БД для хранения хэшей за вами

Выбрана PosgreSQL.

Для данный задачи подойдёт в принципе любая СУБД.

Для создания таблиц используется инструмент `goose`.

[Директория с миграцией](./migrations/).

---

### К Сервисам:

#### Для межсервисного взаимодействия используйте gRPC

Описание gRPC находится в файле [hashcalc.proto](./hashcalc/api/hashcalc.proto).

Сгенерированный пакет [grpchashcalc](./hashcalc/pkg/grpchashcalc/).

Реализация сервера в пакете [grpchandlers](./hashcalc/internal/grpchandlers/).

Реализация клиента в пакете [grpccalc](./hashkeeper/internal/interfaces/grpccalc/).

#### Выбор способа хранения конфигураций (данные для соединения с бд, адрес соседнего сервиса и т.д.) за вами

Конфигурация в оба сервиса передается через флаги и переменные окружения.

Для описания конфигурационных данных сервиса используется модуль [go-arg](https://github.com/alexflint/go-arg).

Позволяет делать структуры, в которых при помощи тегов для каждого поля описывается какой должен быть флаг, переменная окружения и значение по умолчанию.

```go
type LogCfg struct {
	LogLevel    int    `arg:"--log-level,env:LOG_LEVEL" default:"4" help:"0-panic, 1-fatal, 2-error, 3-warn, 4-info, 5-debug, 6-trace"`
	LogGELF     bool   `arg:"--log-gelf,env:LOG_GELF" default:"false" help:"Enable of disable GELF format of logs"`
	LogURL      string `arg:"--log-url,env:LOG_URL" default:"localhost:12201" help:"Host and port of log server, keep it empty to disable sending logs to server"`
	LogHostname string `arg:"--log-hostname,env:LOG_HOSTNAME" default:"localhost" help:"Name of instance in logs"`
}
```

#### Сервисы должны логгировать свое поведение

Для логирования используется модуль [logrus](https://github.com/sirupsen/logrus).

Сделаны обёртки для настройки вывода логов в `graylog` в пакете [hashlog](./hashkeeper/pkg/hashlog/).

#### Service 2 должен реализовывать спецификацию,описанную в swagger’

Описание REST API находится в файле [hashkeeper.yaml](./hashkeeper/pkg/interfaces/restapi/spec/hashkeeper.yaml).

Для генерации используется докер образ `quay.io/goswagger/swagger`, который завернут в скрипт [swagger.sh](./etc/scripts/swagger.sh).

Конкретные аргументы запуска можно посмотреть в файле [Makefile](./Makefile), цель `swagger-gen-srv`.

Сгенерированный сервер находится [в этой директории](./hashkeeper/internal/interfaces/restapi/).

Реализация ручек расположена [в этом файле](./hashkeeper/internal/interfaces/restapi/server/configure_hashkeeper.go).

####  Service 1 подсчет хэшей должен быть организован параллельно

Подсчёт хешей реализован в пакете [strhash](./hashcalc/internal/strhash/).

Для параллельного вычисления используются каналы для подачи и получения данных и отдельные рутины для каждой строки.

Исходные строки подаются через gRPC stream и вычисления хешей начинаются ещё до закрытия стрима.

Результаты так же выдаются через стрим, для другого сервиса.

#### Service 1 должен быть покрыт тестами

Unit тесты сделаны для пакета [strhash](./hashcalc/internal/strhash/).

#### Оба сервиса должны находиться в контейнерах и подниматься в связке с помощью docker-compose

dockerfile для каждого сервиса и общий docker-compose файл расположены в директории [deployments](./deployments/).

#### Оба сервиса должны корректно обрабатывать сигналы о завершении работы (graceful shutdown)

В сервисе `hashcalc` штатное завершение обрабатывается в [главной функции](./hashcalc/cmd/main.go).

```go
	<-sigTerm
	// Stop
	log.Info("hashcalc is shutting down ...")
	grpcServer.Stop()
```

В сервисе `hashkeeper` штатное завершение привязано к реализации уже [сгенерированного сервера](./hashkeeper/internal/interfaces/restapi/server/configure_hashkeeper.go).

```go
	api.PreServerShutdown = func() {
		_server.app.Shutdown()
	}
```

### К логам:

#### Сервисы должны писать логи в формате GELF

Для включения вывода в GELF формате в настройках переменных окружения стоит `LOG_GELF: "true"` в [docker-compose](./deployments/docker-compose.yml).

#### Логирование должно быть сквозным

Сквозное логирование запроса привязано к контексту запроса и следующих за ним операциях.

В обёртке [hashlog](./hashkeeper/pkg/hashlog/) реализованы добавление получения `requestID` к контексту.

Идентификатор запроса также передаётся от `hashkeeper` к `hashcalc` через gRPC посредством [interceptor](./hashcalc/pkg/grpchashcalc/request_id_interceptors.go).

#### К логам должен быть добавлено поле со Stack Trace’ом. Для получения трейса можно использовать:

В пакете [hashlog](./hashkeeper/pkg/hashlog/error.go) реализованы добавление stack trace к ошибкам и последующему его выводу в лог.

#### Продумайте использование уровней логирования (вам надо будет аргументировать во время презентации ваше решение)

- Info - информация о запуске и завершении
- Wark - некоторые сообщения о запуске отладочных инструментов
- Error - вывод ошибок в запросах и функции Main()
- Fatal - вывод из завершение работы приложения
- Debug - информация по запуску и завершению запроса
- Trace - содержимое полей запроса и ответа
