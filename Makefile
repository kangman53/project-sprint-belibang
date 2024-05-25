build-dev:
	docker compose build

restart-dev:
	docker restart belibang-web

run-dev:
	docker compose up -d

logs-web:
	docker logs -f --tail 100 belibang-web

logs-db:
	docker logs -f belibang-db

check-db:
	docker exec -it belibang-db psql -U belibang -d belibang-db

clear-db:
	docker rm -f -v belibang-db

migrate-db:
	migrate -database "postgres://belibang:password@localhost:5432/belibang-db?sslmode=disable" -path database/migrations up
	
migrate-db-down:
	migrate -database "postgres://belibang:password@localhost:5432/belibang-db?sslmode=disable" -path database/migrations down -all
	
build-prod-linux:
	GOOS=linux GOARCH=amd64 go build -o build/belibang

build-prod-win:
	GOOS=windows GOARCH=amd64 go build -o build/belibang.exe

build-prod-mac:
	GOOS=darwin GOARCH=amd64 go build -o build/belibang

build-prod-docker:
	docker build . -t belibang
	docker tag belibang:latest rereasdev/belibang:latest

docker-push:
	docker push rereasdev/belibang:latest
