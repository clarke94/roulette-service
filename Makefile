.PHONY: help
help: ## Shows all make targets available
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/##//'

.PHONY: serve
serve: ## Serve the app locally
	@docker-compose up --build

.PHONY: test
test: ## Test app
	@go test -race ./...

.PHONY: cover
cover: ## Cover app
	@go test -coverprofile=cover.out ./...

.PHONY: lint
lint: ## Lint app
	@golangci-lint run