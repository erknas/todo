## Описание проекта

Проект представляет собой простое Todo приложение, которое позволяет пользователям создавать, читать, обновлять и удалять задачи.


## Запуск проекта

### Локально
В терминале выполнить команду `make run` и перейти по адресу [http://localhost:7540](http://localhost:7540). Аутентификация осуществляется по паролю, указанному в `.env`

### Dockerfile
В терминале выполнить команды `docker build -t todo .` и `docker run -d -p 7540:7540 -v $(pwd)/scheduler.db:/app/scheduler.db todo`, перейти по адресу [http://localhost:7540](http://localhost:7540). Аутентификация осуществляется по паролю, указанному в `.env`