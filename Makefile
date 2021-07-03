build:
	go build -o ragnarok ./cmd/ragnarok

fmt:
	go fmt ./cmd/*
	go fmt ./internal/*
