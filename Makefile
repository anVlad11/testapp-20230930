#!/usr/bin/make

-include include.make

run-local:
	go run ./cmd/app/ --config-path ./config.yaml