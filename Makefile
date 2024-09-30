
build:
	docker build -t capymeme-bot .

up:
	docker-compose up -d

down:
	docker-compose down

delete:
	docker rmi capymeme-bot --force

.PHONY: build up down delete