

# (RU) CRUD API Библиотека для STEAM игр 

# Содержание

- [Описание](#Описание)
- [Список пакетов](#Список-пакетов)
- [Установка](#Установка)
- [Использование](#Использование)
- [Настройка](#Настройка)
- [Логи](#Логи)


# Описание
Общение с API происходит непосредствеено через JSON.
Основная информацию об игре берется из SteamAPI, дополнительную информацию пользователь заполняет сам.

* Реализована аунтификация пользователя, хеширование паролей. Токен сессии выдается на 2 дня.
* Реализовано логирование ошибок, общей информации.
* Приложение докерезировано.
* Реализовано кэширование запроса на вывод количества сущностей с использованием Redis.
* Библиотека, пользователи и сессии хранятся в PostgreSQL.

Таблицы PostgreSQL:
- games
- developers
- publishers
- sessions
- users

Приложение разработано в рамках выпускной квалификационной работы.
Основной задачей было изучение технологий применяемых в бэкенд разработке. 

# Список пакетов

1. [zap](https://github.com/uber-go/zap) - Логер 
2. [random-data](https://github.com/Pallinder/go-randomdata) - Генератор рандомных значений
3. [validator](https://github.com/go-playground/validator/v10) - Валидатор 
4. [pq](https://github.com/lib/pq) - Драйвер для работы с PostgreSQL
5. [redis](https://github.com/go-redis/redis/v8) - Драйвер для работы с Redis

# Установка

### Скачать Docker 

**Linux** 
```$xslt
https://docs.docker.com/engine/install/
```
**Windows**
```$xslt
https://docs.docker.com/desktop/install/windows-install/
```
**MacOS**
```$xslt
https://docs.docker.com/desktop/install/mac-install/
```

### Скопировать репозиторий, зайти в корень
```$xslt
git clone https://github.com/0mwa/steamlib-crud-api.git
```
```$xslt
cd steamlib-crud-api
```

### Запустить docker-compose
```$xslt
docker-compose up -d
```

### Провести миграцию для добавления необходимых таблиц в PostgreSQL
```$xslt
sh initDataBase.sh
```

# Использование

### Создание пользователя
```$xslt
curl -X POST --data '{"login":"user", "password":"user"}' http://localhost:3333/register
```

### Аутентификация
```$xslt
curl -X POST --data '{"login":"admin_1", "password":"admin_1"}' http://localhost:3333/auth
```

### Добавление сущностей
```$xslt
curl -X POST -H "token: 6271879382257613538Matthew" http://localhost:3333/games/add/440
```
```$xslt
curl -X POST -H "token: 6271879382257613538Matthew" http://localhost:3333/developres/add/440
```
```$xslt
curl -X POST -H "token: 6271879382257613538Matthew" http://localhost:3333/publishers/add/440
```

### Получение списка всех сущностей 
```$xslt
curl -X GET -H "token: 6271879382257613538Matthew" http://localhost:3333/games
```

### Получение конкретной сущности 
```$xslt
curl -X GET -H "token: 6271879382257613538Matthew" http://localhost:3333/games/440
```

### Получение конкретной сущности
```$xslt
curl -X GET -H "token: 6271879382257613538Matthew" http://localhost:3333/games/440
```

### Удаление сущнности
```$xslt
curl -X DELETE -H "token: 6271879382257613538Matthew" http://localhost:3333/games/delete/440
```

### Дополнение информации о сущности
```$xslt
curl -X PUT -H "token: 6271879382257613538Matthew" http://localhost:3333/games/update/440 --data '{
                                                           "name": "Team Fortress 2",
                                                           "img": "img.png",
                                                           "description": "",
                                                           "rating": 10,
                                                           "developer_id": 4,
                                                           "publisher_id": 4,
                                                           "steam_id": 440
                                                          }'
```

### Получение количества всех сущностей 
```$xslt
curl -X GET -H "token: 6271879382257613538Matthew" http://localhost:3333/get_counter
```

# Настройка

Пароль, хост, юзер, порт.
При желании настраиваются в .env.

# Логи

### Чтобы посмотреть логи, нужно узнать наименование контейнера приложения.
```$xslt
docker ps 
```
```$xslt
docker logs app-name
```