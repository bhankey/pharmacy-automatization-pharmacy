.PHONY:

build:
	go mod download && go build -o ./server ./cmd/server/main.go

docker:
	docker-compose -f .build/docker-compose.yaml up

docker-down:
	docker-compose -f .build/docker-compose.yaml down

migrations-down:
	MIGRATIONS_STATUS=down docker-compose -f docker-compose-down-migrations.yaml up

run: .build
	./server

test:
	go test -v ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

proto:
	protoc --go_out="./pkg/" --go_opt=paths=import \
			--go-grpc_out="./pkg/" --go-grpc_opt=paths=import \
 			--validate_out="lang=go:./pkg" \
 			-I "${GOPATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.7" -I api/grpc api/grpc/*.proto
