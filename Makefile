.SILENT:

# Создать образ для сервиса основного приложения.
buildimg:
	docker build -t appointment-api ./appointment
#	docker build -t dp210go_auth    ./auth
	docker build -t doctor-api      ./doctor
	docker build -t user-api        ./user
up:
	docker-compose up -d
stop:
	docker-compose stop
#cleanjunk:
#	docker system prune


#logsall:
#	docker-compose logs -f
#
#logswebapp:
#	docker-compose logs -f webapp
#
#logsdb:
#	docker-compose logs -f db
#
#logsredis:
#	docker-compose logs -f redis

upall:buildimg up #cleanjunk

