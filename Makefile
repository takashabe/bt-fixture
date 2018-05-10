SUBPACKAGES := $(shell go list ./...)

.DEFAULT_GOAL := help

##### Development

.PHONY: deps test vet lint

deps: ## Setup dependencies package
	dep ensure

test: ## Run go test
	go test -v -p 1 $(SUBPACKAGES)

vet: ## Check go vet
	go vet $(SUBPACKAGES)

lint: ## Check golint
	golint $(SUBPACKAGES)

##### Utilities

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
