all: generate

run:
	go run cmd/gRPC_service/main.go --config=./config/local.yaml

migrations:
	go run ./cmd/migrator/main.go --storage-path=./storage/gRPC_service.db --migrations-path=./migrations

migrations-test:
	go run ./cmd/migrator/main.go --storage-path=./storage/gRPC_service.db --migrations-path=./tests/migrations --migrations-table=migrations_test