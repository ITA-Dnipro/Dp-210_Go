.SILENT:

build:
	go build -o ./bin/dp210goapp ./

buildimg:
	docker build -t dp210goimg ./

up:
	docker-compose up -d

stop:
	docker-compose stop

logsall:
	docker-compose logs -f

logswebapp:
	docker-compose logs -f webapp

logsdb:
	docker-compose logs -f db

logsredis:
	docker-compose logs -f redis

# UPALL --------------------------------------------------|

buildandbuildimg:build
	docker build -t dp210goimg ./

upall:buildandbuildimg
	docker-compose up -d
