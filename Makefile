all: generate

migrations:
	go run ./cmd/migrator --storage-path=./storage/gRPC_service.db --migrations-path=./migrations 