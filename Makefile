.PHONY develop:
up:
	docker-compose -f docker-compose-dev.yml up -d
build-up:
	docker-compose -f docker-compose-dev.yml up --build -d
build-front:
	docker-compose -f docker-compose-dev.yml up -d --build frontend
build-back:
	docker-compose -f docker-compose-dev.yml up -d --build backend
down:
	docker-compose -f docker-compose-dev.yml down
down-v:
	docker-compose -f docker-compose-dev.yml down -v
