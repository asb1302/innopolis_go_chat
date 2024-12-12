# innopolis_go_chat

Серверная часть приложения для взаимодействия чатов через WebSocket

## Технические подробности

* Адрес прода wss://innopolisgochat-production.up.railway.app
* [Описание задания](https://github.com/asb1302/innopolis_go_chat_client/blob/dev/task%2FREADME.md)
* Для взаимодействия сервера и клиента используется специально созданный для этого go-модуль
  [innopolis_go_chat/tree/asb1302-dev/pkg/chatdata](https://github.com/asb1302/innopolis_go_chat/tree/asb1302-dev/pkg/chatdata)

## Локальная работа с проектом

### Требования

- Docker
- Docker Compose
- Make

### Основные команды Makefile

#### Сборка и запуск проекта

Для сборки и запуска проекта выполните следующую команду:

```sh
make init
```

#### Очистка проекта

Для остановки и удаления контейнеров, образов, томов и зависших контейнеров выполните следующую команду:

```sh
make clean
```

#### Примечание

Убедитесь, что вы настроили переменные окружения(.env) перед запуском команд.