include .env

docker.create.container:
	@echo "creating db container"
	docker run --name ${DB_DOCKER_CONTAINER} -e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} -e MYSQL_DATABASE=${DB_NAME} -p 3306:3306 -d mysql:8 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

docker.start.container:
	@echo "running docker images"
	docker start ${DB_DOCKER_CONTAINER}

show.tables:
	@echo "connecting to my db..."
	docker exec -i ${DB_DOCKER_CONTAINER} mysql -u${DB_USER} -p${MYSQL_ROOT_PASSWORD} ${DB_NAME} -Bse "show tables;"

init.up:
	@echo "initiating tables..."
	docker exec -i ${DB_DOCKER_CONTAINER} mysql -u${DB_USER} -p${MYSQL_ROOT_PASSWORD} ${DB_NAME} < ./migrations/init.up.sql

init.down:
	@echo "Deleting initial creation..."
	docker exec -i ${DB_DOCKER_CONTAINER} mysql -u${DB_USER} -p${MYSQL_ROOT_PASSWORD} ${DB_NAME} < ./migrations/init.down.sql

seed.up:
	@echo "inserting meetings rows"
	docker exec -i ${DB_DOCKER_CONTAINER} mysql -u${DB_USER} -p${MYSQL_ROOT_PASSWORD} ${DB_NAME} < ./migrations/seed.up.sql

seed.down:
	@echo "deleting migration \"seed\""
	docker exec -i ${DB_DOCKER_CONTAINER} mysql -u${DB_USER} -p${MYSQL_ROOT_PASSWORD} ${DB_NAME} < ./migrations/seed.down.sql

# build all binaries
build: clean build.client build.server
	@printf "All binaries build!"

clean:
	@echo "Cleaning..."
	@- rm -f dist/*
	@go clean
	@echo "Cleaned!"

build.client:
	@echo "Building clientend..."
	@go build -o dist/gostripe ./cmd/web
	echo "Front is build!"

build.server:
	@echo "Building serverend..."
	@env go build -o dist/gostripe_api ./cmd/api
	@echo "Backend is built!"

start: start.client start.server
	@echo "Starting the client end..."

start.client:
	@echo "Starting the clientend..."
	@env STRIPE_SECRET=${STRIPE_SECRET} STRIPE_KEY=${STRIPE_KEY} DSN=${DSN} ./dist/gostripe -port=${GOSTRIPE_PORT} &
	@echo "Front end running!"

start.server:
	@echo "Starting the server end..."
	@env USERNAME=${USERNAME} PASSWORD=${PASSWORD} STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} DSN=${DSN} ./dist/gostripe_api -port=${API_PORT} &
	@echo "Back end running!"

restart.server: stop.server build.server start.server

stop: stop.client stop.server
	@echo "All applications stopped"

stop.client:
	@echo "Stopping the client end..."
	@-pkill -SIGTERM -f "gostripe -port=${GOSTRIPE_PORT}"
	@echo "Stopped client end"

stop.server:
	@echo "Stopping the server end..."
	@-pkill -SIGTERM -f "gostripe_api -port=${API_PORT}"
	@echo "Stopped server end"

run.client: stop.client build.client start.client

run.server: docker.start.container stop.server build.server start.server

