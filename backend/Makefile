include .env
export

up:
	@docker compose up -d --build polytweet db-postgres
	@docker compose run --rm db-migrations -path ./migrations -database ${DATABASE_URL} up

down:
	@docker compose down --remove-orphans

migrate-up:
	@docker compose run --rm db-migrations -path ./migrations -database ${DATABASE_URL} up 1

migrate-down:
	@docker compose run --rm db-migrations -path ./migrations -database ${DATABASE_URL} down 1
