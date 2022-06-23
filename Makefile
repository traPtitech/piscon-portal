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
.PHONY: deploy
deploy:
	docker compose -f docker-compose-prod.yml up -d
.PHONY: deploy-build
deploy-build:
	docker compose -f docker-compose-prod.yml up -d --build
