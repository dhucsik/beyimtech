build:
	go build -o ./bin/service ./cmd/main.go

migrate: 
	go install github.com/SergeiSkv/goose/v3/cmd/goose@v3.20.8
	@echo "\033[0;32mMigrating..."
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=postgres://postgres:dias2003dewex@localhost:5432/diploma goose -dir=${PWD}/migrations up

docs:
	swag init -d ./internal/transport/http -g server.go --pd --parseDepth 2 -o ./swagger/docs

lint:
	golangci-lint run