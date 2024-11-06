start:
	docker-compose -f docker.postgres.yaml up --build -d

stop:
	docker-compose -f docker.postgres.yaml down

generate_sql:
	sqlc generate

migrate_up:
	goose -dir sql/schema postgres "postgres://postgres:postgres@127.0.0.1:5432/gator" up

migrate_down:
	goose -dir sql/schema postgres "postgres://postgres:postgres@127.0.0.1:5432/gator" down

db_login:
	docker exec -it gator-db-1 psql -U postgres
