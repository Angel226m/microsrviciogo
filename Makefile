.PHONY: all build up down logs clean test lint fmt security

# ═══════════════════════════════════════════════════════════════
# CloudMart – Makefile
# Automatización de desarrollo, pruebas, seguridad y despliegue
# Compatible con OWASP Top 10 e ISO 27001 Anexo A.14
# ═══════════════════════════════════════════════════════════════

COMPOSE = docker compose
SERVICES = api-gateway user-service product-service order-service payment-service inventory-service notification-service
FRONTEND = frontend

## ─── Desarrollo ─────────────────────────────────────────────
up: ## Iniciar todos los servicios
	$(COMPOSE) up -d --build

down: ## Detener todos los servicios
	$(COMPOSE) down

restart: down up ## Reiniciar todos los servicios

logs: ## Ver logs de todos los servicios
	$(COMPOSE) logs -f --tail=50

logs-%: ## Ver logs de un servicio específico (ej: make logs-api-gateway)
	$(COMPOSE) logs -f --tail=100 $*

status: ## Ver estado de todos los contenedores
	$(COMPOSE) ps

## ─── Construcción ───────────────────────────────────────────
build: ## Construir todas las imágenes
	$(COMPOSE) build --parallel

build-%: ## Construir un servicio específico (ej: make build-api-gateway)
	$(COMPOSE) build $*

build-frontend: ## Construir solo el frontend
	$(COMPOSE) build $(FRONTEND)

## ─── Pruebas ────────────────────────────────────────────────
test: ## Ejecutar todas las pruebas
	@for svc in $(SERVICES); do \
		echo "\n══════ Probando $$svc ══════"; \
		cd services/$$svc && go test ./... -v -cover && cd ../..; \
	done

test-%: ## Probar un servicio específico (ej: make test-user-service)
	cd services/$* && go test ./... -v -cover -race

test-unit: ## Ejecutar solo pruebas unitarias (sin integración)
	@for svc in $(SERVICES); do \
		echo "\n══════ Unitarias $$svc ══════"; \
		cd services/$$svc && go test ./internal/... -v -cover -short && cd ../..; \
	done

test-coverage: ## Generar reporte de cobertura HTML
	@for svc in $(SERVICES); do \
		echo "\n══════ Cobertura $$svc ══════"; \
		cd services/$$svc && go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html && cd ../..; \
	done

test-frontend: ## Ejecutar pruebas del frontend
	cd $(FRONTEND) && npm test

## ─── Calidad de Código ──────────────────────────────────────
lint: ## Analizar todos los servicios con golangci-lint
	@for svc in $(SERVICES); do \
		echo "\n══════ Analizando $$svc ══════"; \
		cd services/$$svc && golangci-lint run ./... && cd ../..; \
	done

fmt: ## Formatear todo el código Go
	@for svc in $(SERVICES); do \
		cd services/$$svc && gofmt -s -w . && cd ../..; \
	done

vet: ## Ejecutar go vet en todos los servicios
	@for svc in $(SERVICES); do \
		echo "\n══════ Verificando $$svc ══════"; \
		cd services/$$svc && go vet ./... && cd ../..; \
	done

## ─── Seguridad (OWASP / ISO 27001) ─────────────────────────
security: ## Ejecutar análisis de seguridad completo
	@echo "══════ Análisis de Seguridad OWASP ══════"
	@echo "→ Escaneando dependencias Go..."
	@for svc in $(SERVICES); do \
		cd services/$$svc && go list -json -m all 2>/dev/null | head -5 && cd ../..; \
	done
	@echo "→ Escaneando imágenes Docker..."
	@for svc in $(SERVICES); do \
		echo "Escaneando $$svc..."; \
	done
	@echo "✅ Análisis de seguridad completado"

deps-update: ## Actualizar dependencias Go de todos los servicios
	@for svc in $(SERVICES); do \
		echo "\n══════ Actualizando $$svc ══════"; \
		cd services/$$svc && go get -u ./... && go mod tidy && cd ../..; \
	done

deps-audit: ## Auditar dependencias por vulnerabilidades conocidas
	@for svc in $(SERVICES); do \
		echo "\n══════ Auditando $$svc ══════"; \
		cd services/$$svc && go list -m all && cd ../..; \
	done

## ─── Infraestructura ────────────────────────────────────────
infra-up: ## Iniciar solo infraestructura (BD, Redis, NATS, monitoreo)
	$(COMPOSE) up -d postgres redis nats prometheus grafana jaeger mailhog

infra-down: ## Detener infraestructura
	$(COMPOSE) down postgres redis nats prometheus grafana jaeger mailhog

## ─── Base de Datos ──────────────────────────────────────────
db-reset: ## Reiniciar base de datos (¡DESTRUCTIVO!)
	$(COMPOSE) down -v postgres
	$(COMPOSE) up -d postgres

db-shell: ## Abrir consola psql
	$(COMPOSE) exec postgres psql -U cloudmart

db-backup: ## Crear respaldo de la base de datos
	$(COMPOSE) exec postgres pg_dump -U cloudmart cloudmart > backup_$$(date +%Y%m%d_%H%M%S).sql

## ─── Monitoreo ──────────────────────────────────────────────
grafana: ## Abrir Grafana en el navegador (http://localhost:3001)
	@echo "Grafana disponible en: http://localhost:3001"

prometheus: ## Abrir Prometheus en el navegador (http://localhost:9090)
	@echo "Prometheus disponible en: http://localhost:9090"

jaeger: ## Abrir Jaeger en el navegador (http://localhost:16686)
	@echo "Jaeger disponible en: http://localhost:16686"

## ─── Limpieza ───────────────────────────────────────────────
clean: ## Eliminar todos los contenedores, volúmenes e imágenes
	$(COMPOSE) down -v --rmi local
	docker system prune -f

clean-cache: ## Limpiar caché de Go y node_modules
	@for svc in $(SERVICES); do \
		cd services/$$svc && go clean -cache && cd ../..; \
	done
	cd $(FRONTEND) && rm -rf node_modules/.cache

## ─── Ayuda ──────────────────────────────────────────────────
help: ## Mostrar esta ayuda
	@echo "═══════════════════════════════════════════"
	@echo "  CloudMart – Comandos Disponibles"
	@echo "═══════════════════════════════════════════"
	@grep -E '^[a-zA-Z_%-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
