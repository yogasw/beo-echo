.PHONY: dev server frontend install-air help

help: ## Show this help menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-air: ## Install air-verse/air for Go hot reloading
	go install github.com/air-verse/air@latest

server: ## Run the Go backend using air
	cd backend && $$(go env GOPATH)/bin/air

frontend: ## Run the Svelte frontend
	cd frontend && npm install && npm run dev

dev: ## Run both backend and frontend concurrently
	$(MAKE) -j2 server frontend
