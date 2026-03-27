.PHONY: all build up down logs clean test lint proto

# ═══════════════════════════════════════════════════════════════
# CloudMart – Makefile
# ═══════════════════════════════════════════════════════════════

COMPOSE = docker compose
SERVICES = api-gateway user-service product-service order-service payment-service inventory-service notification-service

## ─── Development ────────────────────────────────────────────
up: ## Start all services
	$(COMPOSE) up -d --build

down: ## Stop all services
	$(COMPOSE) down

restart: down up ## Restart all services

logs: ## Tail logs for all services
	$(COMPOSE) logs -f --tail=50

logs-%: ## Tail logs for a specific service (e.g., make logs-api-gateway)
	$(COMPOSE) logs -f --tail=100 $*

## ─── Build ──────────────────────────────────────────────────
build: ## Build all service images
	$(COMPOSE) build --parallel

build-%: ## Build a specific service (e.g., make build-api-gateway)
	$(COMPOSE) build $*

## ─── Testing ────────────────────────────────────────────────
test: ## Run all tests
	@for svc in $(SERVICES); do \
		echo "\n══════ Testing $$svc ══════"; \
		cd services/$$svc && go test ./... -v -cover && cd ../..; \
	done

test-%: ## Test a specific service (e.g., make test-user-service)
	cd services/$* && go test ./... -v -cover -race

## ─── Code Quality ───────────────────────────────────────────
lint: ## Lint all services
	@for svc in $(SERVICES); do \
		echo "\n══════ Linting $$svc ══════"; \
		cd services/$$svc && golangci-lint run ./... && cd ../..; \
	done

fmt: ## Format all Go code
	@for svc in $(SERVICES); do \
		cd services/$$svc && gofmt -s -w . && cd ../..; \
	done

## ─── Infrastructure ─────────────────────────────────────────
infra-up: ## Start only infrastructure (DB, Redis, NATS, monitoring)
	$(COMPOSE) up -d postgres redis nats prometheus grafana jaeger mailhog

infra-down: ## Stop infrastructure
	$(COMPOSE) down postgres redis nats prometheus grafana jaeger mailhog

## ─── Database ───────────────────────────────────────────────
db-reset: ## Reset database
	$(COMPOSE) down -v postgres
	$(COMPOSE) up -d postgres

db-shell: ## Open psql shell
	$(COMPOSE) exec postgres psql -U cloudmart

## ─── Cleanup ────────────────────────────────────────────────
clean: ## Remove all containers, volumes, and images
	$(COMPOSE) down -v --rmi local
	docker system prune -f

## ─── Help ───────────────────────────────────────────────────
help: ## Show this help
	@grep -E '^[a-zA-Z_%-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
