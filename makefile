dbUrl="postgres://postgres:postgres@localhost:5432/site"
path="./migrations"

build: 
	go build -o site

migrate-up:
	migrate -database=${dbUrl} -path=${path} up

migrate-down:
	migrate -database=${dbUrl} -path=${path} down

migrate-seed:
	psql -U postgres -h localhost -f seed.sql -d site

migrate-reset: migrate-down migrate-up migrate-seed
