lint:
	golangci-lint run -c ./golangci.yml ./...

build:
	go build -o bin/app ./cmd/rssaggregator

run: build
	./bin/app --config=./config/local.yml

test:
	go test -v ./... -count=1

generate-mocks:
	mockery

clean-mocks:
	rm -rf ./internal/mocks/*

swagger:
	swag init --dir ./cmd/rssaggregator,./api --output ./docs --parseDependency --parseInternal

