test_files	= $(shell go list ./... )

.PHONY: mock
mock:
	bin/mock-generator.sh

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test -v -race $(test_files)

.PHONY: test-coverage
test-coverage:
	go test -v -race $(test_files) -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out

.PHONY: docker-compose
docker-compose:
	docker-compose up --build