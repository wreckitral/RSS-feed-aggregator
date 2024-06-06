build:
	@go build -o bin/RSS-feed-aggregator

run: build
	@./bin/RSS-feed-aggregator

test:
	@go test ./... -v
