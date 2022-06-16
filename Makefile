.PHONY: *

run:
	go run . -v cron
	go run . -v -- snapshots

test:
	go test ./... -v

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

build:
	goreleaser release --rm-dist --skip-publish --snapshot

release:
	goreleaser release --rm-dist
