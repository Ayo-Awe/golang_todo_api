include .env

dev:
	air

migrate-up:
	migrate -path migrations -verbose -database ${DATABASE_URL} up

migrate-down:
	migrate -path migrations -verbose -database ${DATABASE_URL} down

migrate-force:
	migrate -path migrations -database ${DATABASE_URL} -verbose force ${version}

new-migration:
	migrate create -dir migrations -seq -ext .sql ${name}

migrate-drop:
	migrate -path migrations -database ${DATABASE_URL} -verbose drop

.PHONY: migrate-up migrate-down new-migration migrate-force dev
