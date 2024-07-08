all: generate

run:
	go run cmd/gRPC_service/main.go --config=./config/local.yaml

migrations:
	go run ./cmd/migrator --storage-path=./storage/gRPC_service.db --migrations-path=./migrations 