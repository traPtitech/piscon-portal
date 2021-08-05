.PHONY docker:
up:
	docker-compose -f docker-compose_dev.yml up -d
build-up:
	docker-compose -f docker-compose_dev.yml up --build -d
build-front:
	docker-compose -f docker-compose_dev.yml up --build frontend -d
build-back:
	docker-compose -f docker-compose_dev.yml up --build backend -d
down:
	docker-compose -f docker-compose_dev.yml down