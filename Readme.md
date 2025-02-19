# File Service

Клиент - серверное приложение для работы с файлмаи

    Функции клиента:    

    ПОЛУЧИТЬ список файлов
    ПОЛУЧИТЬ файл по его имени с сервера и сохранить на клиенте
    ОТПРАВИТЬ файл на сервер
    ОБНОВИТЬ файл на сервере
    УДАЛИТЬ файл с сервера


## 
Для запуска:

    cd server/

    docker build -t fileservice-server .

    cd client/

    docker build -t fileservice-client .

    docker run fileservice-server

    docker run --rm -it fileservice-client

## Запуск сервера:
    make server

## Запуск клиента:
    make client
