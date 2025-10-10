migrate-up:
	@set -a; source .env; set +a; \
	go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres \
	"postgres://$$DB_USER:$$DB_PASSWORD@localhost:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL" up

migrate-down:
	@set -a; source .env; set +a; \
	go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres \
	"postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL" down

run:
	@docker compose up -d
	@echo "Waiting for Postgres healthy..."
	@until [ "`docker inspect -f {{.State.Health.Status}} dev_postgres`" = "healthy" ]; do \
		sleep 1; \
	done
	@$(MAKE) migrate-up
	@docker compose up

logs:
	@docker compose logs -f
