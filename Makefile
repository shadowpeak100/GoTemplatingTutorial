#!/usr/bin/make -f

run:
	go run templatePresentation/main.go

compile:
	go build templatePresentation/main.go

.PHONY: run compile