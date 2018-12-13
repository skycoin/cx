.DEFAULT_GOAL := help
.PHONY: build-parser build install-deps install test update-golden-files

build-parser: ## Generate lexer and parser for CX grammar
	# TODO: Implement

build: install-build-deps build-parser ## Build CX from sources
	# TODO: Implement

install-build-deps:

install: build ## Install CX from sources

test: build ## Run CX test suite.

update-golden-files: ## Update golden files used in CX test suite

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

