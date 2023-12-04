include .env

init.up:
	@echo "initiating tables..."
	@env mysql -uroot -p'${PASSWORD}' widgets < ./migrations/init.up.sql

show-tables:
	@echo "connecting to my db widgets"
	@env mysql -uroot -p'${PASSWORD}' widgets -Bse "show tables;"


init.down:
	@echo "Deleting initial creation..."
	@env mysql -uroot -p'${PASSWORD}' widgets < ./migrations/init.down.sql

seed.up:
	@echo "inserting session row"
	@env mysql -uroot -p'${PASSWORD}' widgets < ./migrations/seed.up.sql

seed.down:
	@echo "deleting migration \"seed\""
	@env mysql -uroot -p'${PASSWORD}' widgets < ./migrations/seed.down.sql


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
	@go build -o dist/gostripe_api ./cmd/api
	@echo "Backend is built!"

start: start.client start.server
	@echo "Starting the client end..."

start.client:
	@echo "Starting the clientend..."
	@env STRIPE_SECRET=${STRIPE_SECRET} STRIPE_KEY=${STRIPE_KEY} ./dist/gostripe -port=${GOSTRIPE_PORT} &
	@echo "Front end running!"

start.server: build.server
	@echo "Starting the server end..."
	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} ./dist/gostripe_api -port=${API_PORT} &
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

run.server: stop.server build.server start.server

