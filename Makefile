setup-git:
	git config --global url."ssh://git@github.com".insteadOf "https://github.com"

install-dependencies:
	go install go.uber.org/mock/mockgen@latest

migrate-db-up:
	go run ./cmd/migration/main.go

migrate-db-down:
	docker run --rm -v $(shell pwd)/.database/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://user:password@localhost:15432/academy_portal_service?sslmode=disable down 1

autogenerate-code:
	docker run --rm -v $(shell pwd):/src -w /src sqlc/sqlc generate

autogenerate-mock:
	mockgen --package mockdatabase --destination autogenerate/database/mock/store.go github.com/WhiteMaks/go_academy_portal_service/autogenerate/database Store

up-test-environment:
	docker compose up -d
	sleep 30

test:
	go clean -testcache && go test -cover -v ./...

build-migration:
	go build -o .bin/academy_portal_service/migration ./cmd/migration/main.go

.PHONY: setup-git install-dependencies migrate-db-up migrate-db-down autogenerate-code autogenerate-mock up-test-environment test build-migration