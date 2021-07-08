build:
	go build -o ragnarok ./cmd/ragnarok

fmt:
	go fmt ./cmd/*
	go fmt ./internal/*

test:
	go test ./cmd/*
	go test ./internal/*

etags:
	rm -f TAGS
	find . -name "*.go" -print | etags -
