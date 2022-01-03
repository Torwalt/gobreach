# GoBreach

Golang service to retrieve, cache and serve account breach data from different backends.

This is a WIP and a test project to explore golang.

## Prerequisites

- golang
- docker
- make

## Init

- cp example.env .env
- fill out missing env vars
- `make run-docker-composed`
- `make init-db`

## Running Tests

- `make test`
- run specific integration test:
  - `go test ./... '-run=^TestBreachDBCreateIntegration$'`

