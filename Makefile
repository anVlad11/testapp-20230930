#!/usr/bin/make

-include include.make

run-local:
	go run ./cmd/app/ --root https://meduza.io