# development
.PHONY: up
up:
	docker compose -p piscon-portal-dev -f docker-compose-dev.yml up -d
.PHONY: build-up
build-up:
	docker compose -p piscon-portal-dev -f docker-compose-dev.yml up -d --build
.PHONY: build-front
build-front:
	docker compose -p piscon-portal-dev -f docker-compose-dev.yml up -d --build frontend
.PHONY: build-back
build-back:
	docker compose -p piscon-portal-dev -f docker-compose-dev.yml up -d --build backend
.PHONY: down
down:
	docker compose -p piscon-portal-dev -f docker-compose-dev.yml down
.PHONY: down-v
down-v:
	docker compose -p piscon-portal-dev -f docker-compose-dev.yml down -v

# production
.PHONY: deploy
deploy:
	docker compose -f docker-compose-prod.yml up -d --build

.PHONY: enter-backend
enter-backend:
	docker container exec -it piscon-portal-backend bash

.PHONY: enter-frontend
enter-frontend:
	docker container exec -it piscon-portal-frontend bash

.PHONY: enter-db
enter-db:
	docker container exec -it piscon-portal-db bash

.PHONY: log-backend
log-backend:
	docker logs piscon-portal-backend -f
