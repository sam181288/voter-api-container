# Purpose: Makefile for voter api
SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "	   build				Build voter api"
	@echo "	   build-run				Build and run voter api"
	@echo "	   run					Run docker compose up"
	@echo "	   stop					Stop running voter API by docker compose down"

.PHONY: run
run:
	docker-compose up

.PHONY: build
build:
	./build-better-docker.sh

.PHONY: build-run
build-run: build run

.PHONY: stop
stop:
	docker-compose down