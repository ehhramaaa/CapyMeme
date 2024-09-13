
build:
	docker build -t capymeme .

up:
	docker-compose up -d

down:
	docker-compose down

delete:
	docker rmi capymeme --force

.PHONY: build up down delete