.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o bin/app main.go

.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: clean
clean:
	rm -rf bin/