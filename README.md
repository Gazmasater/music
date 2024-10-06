# Music API

Это API для управления музыкальной библиотекой, позволяющее получать, добавлять, обновлять и удалять песни.



## Установка

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/Gazmasater/music

2. Запустите контейнер: Используйте Docker Compose для сборки и запуска контейнеров:

docker-compose up --build

3. Проверка работы: После запуска контейнера  API будет доступно по адресу http://localhost:8081.

4. Запуск сваггера в браузере

http://localhost:8081/swagger/index.html

5. Остановка контейнера: Чтобы остановить запущенные контейнеры, выполните:

docker-compose down

## Как форматировать песню для использования в данном api

1. cd cmd

2. Перенести текст песни с разделенными куплетами в cmd/input.txt

3. go run .

4. Отформатированный текст будет в cmd/output.txt

5. Открыть в сваггере "Изменение данных песен" . В теле запроса внести полученную строку в файле output.txt
вместо   "string" в поле    "verses"

{
  "artist_name": "Исполнитель",
  "group_link": "http://example.com",
  "release_date": "1985-02-05",
  "song_name": "Название песни",
  "text": {
    "verses": [
      "string"
    ]
  }
}


