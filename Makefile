SERVICE_NAME=innopolis_go_chat
PORT=18002

.PHONY: default
default: build up

.PHONY: build
build:
	docker-compose build

# Запуск контейнеров в фоне
.PHONY: up
up:
	docker-compose up -d

.PHONY: init
init: build up

.PHONY: exec
exec:
	@container_id=$$(docker ps -q --filter "name=$(SERVICE_NAME)"); \
	if [ -n "$$container_id" ]; then \
		docker exec -it $$container_id sh; \
	else \
		echo "Service $(SERVICE_NAME) is not running."; \
	fi

.PHONY: stop
stop:
	docker-compose stop

.PHONY: clean
clean:
	docker-compose down --rmi all --volumes --remove-orphans

.PHONY: check

