.SILENT:

# Создать образ для сервиса основного приложения.
buildimg:
	docker build -t dp210goimg ./

# Поднять все сервисы.
up:
	docker-compose up -d

# Остановить все сервисы.
stop:
	docker-compose stop

# Запуск в терминале логов всех сервисов.
logsall:
	docker-compose logs -f

# Запуск в терминале лога сервиса приложения.
logswebapp:
	docker-compose logs -f webapp

# Запуск в терминале лога сервиса db.
logsdb:
	docker-compose logs -f db

# Запуск в терминале лога сервиса redis.
logsredis:
	docker-compose logs -f redis

# Поднять все с созданием образа для сервиса приложения.
upall:buildimg
	docker-compose up -d
