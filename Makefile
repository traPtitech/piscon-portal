.PHONY develop:
up:
	docker-compose -p piscon-portal-dev -f docker-compose-dev.yml up -d
build-up:
	docker-compose -p piscon-portal-dev -f docker-compose-dev.yml up -d --build
build-front:
	docker-compose -p piscon-portal-dev -f docker-compose-dev.yml up -d --build frontend
build-back:
	docker-compose -p piscon-portal-dev -f docker-compose-dev.yml up -d --build backend
down:
	docker-compose -p piscon-portal-dev -f docker-compose-dev.yml down
down-v:
	docker-compose -p piscon-portal-dev -f docker-compose-dev.yml down -v
