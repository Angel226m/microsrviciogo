# CloudMart – Ecommerce con Service Mesh

```
╔═══════════════════════════════════════════════════════════════════╗
║                        C L O U D M A R T                         ║
║              Production-Grade Ecommerce Platform                  ║
║         Microservices · Istio · Observability · GitOps            ║
╚═══════════════════════════════════════════════════════════════════╝
```

## 🏗️ Architecture

```
                    ┌──────────────┐
                    │   Frontend   │  React 19 · TypeScript · Tailwind 4
                    │   :3000      │
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │ API Gateway  │  Routing · Auth · Rate Limiting
                    │   :8080      │  Circuit Breaker · Load Balancing
                    └──────┬───────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
    ┌─────▼─────┐   ┌─────▼─────┐   ┌─────▼─────┐
    │   User    │   │  Product  │   │   Order   │
    │  Service  │   │  Service  │   │  Service  │
    │   :8081   │   │   :8082   │   │   :8083   │
    └───────────┘   └───────────┘   └─────┬─────┘
                                          │
                    ┌─────────────────────┼──────────────┐
                    │                     │              │
              ┌─────▼─────┐        ┌─────▼─────┐  ┌────▼──────┐
              │  Payment  │        │ Inventory │  │Notification│
              │  Service  │        │  Service  │  │  Service   │
              │   :8084   │        │   :8085   │  │   :8086    │
              └───────────┘        └───────────┘  └────────────┘

    ┌─────────────────────────────────────────────────────────┐
    │                   Infrastructure                         │
    │  PostgreSQL · Redis · NATS JetStream                     │
    │  Prometheus · Grafana · Jaeger · Istio                   │
    └─────────────────────────────────────────────────────────┘
```

## 🛠️ Tech Stack

| Layer          | Technology                                    |
|----------------|-----------------------------------------------|
| Frontend       | React 19.3, TypeScript 5.7, Tailwind CSS 4.2  |
| Backend        | Go 1.22, Hexagonal Architecture               |
| API Gateway    | Custom Go gateway with middleware chain        |
| Messaging      | NATS JetStream (async events)                 |
| Database       | PostgreSQL 16 (schema-per-service)             |
| Cache          | Redis 7 (sessions, product cache)              |
| Service Mesh   | Istio (VirtualServices, DestinationRules)      |
| Observability  | Prometheus + Grafana + Jaeger                  |
| CI/CD          | GitHub Actions + ArgoCD (GitOps)               |
| Containers     | Docker + Docker Compose + Kubernetes           |

## 📁 Project Structure

```
microservicios/
├── docker-compose.yml          # Local orchestration
├── Makefile                    # Developer commands
├── services/
│   ├── api-gateway/            # Reverse proxy + middleware
│   ├── user-service/           # Auth, users, addresses
│   ├── product-service/        # Catalog, categories, reviews
│   ├── order-service/          # Orders, order items
│   ├── payment-service/        # Transactions, payment processing
│   ├── inventory-service/      # Stock, movements, reservations
│   └── notification-service/   # Email, push notifications
├── frontend/                   # React SPA
├── k8s/                        # Kubernetes manifests
│   ├── base/                   # Base deployments
│   ├── istio/                  # Service mesh configs
│   └── monitoring/             # Prometheus, Grafana, Jaeger
├── infra/
│   ├── prometheus/             # Prometheus config
│   ├── grafana/                # Dashboards & provisioning
│   └── argocd/                 # GitOps application manifests
└── .github/workflows/          # CI/CD pipelines
```

## 🚀 Quick Start

```bash
# Clone & start everything
git clone <repo-url>
cd microservicios

# Start all services with Docker Compose
make up

# Or start only infrastructure first
make infra-up

# Check logs
make logs

# Run tests
make test
```

## 🔗 Service URLs (Local)

| Service        | URL                          |
|----------------|------------------------------|
| Frontend       | http://localhost:3000         |
| API Gateway    | http://localhost:8080         |
| Grafana        | http://localhost:3001         |
| Prometheus     | http://localhost:9090         |
| Jaeger UI      | http://localhost:16686        |
| NATS Monitor   | http://localhost:8222         |
| MailHog        | http://localhost:8025         |

## 🔐 Default Credentials

| Service  | User                  | Password          |
|----------|-----------------------|-------------------|
| App      | admin@cloudmart.dev   | admin123          |
| App      | customer@cloudmart.dev| admin123          |
| Grafana  | admin                 | cloudmart         |
| Postgres | cloudmart             | cloudmart_secret  |

## 📊 Hexagonal Architecture (each service)

```
service/
├── cmd/main.go                     # Entry point
├── internal/
│   ├── domain/                     # Business core (0 dependencies)
│   │   ├── model/                  # Entities & value objects
│   │   ├── port/                   # Interfaces (driven & driving)
│   │   └── event/                  # Domain events
│   ├── application/                # Use cases / services
│   │   └── service/
│   └── infrastructure/             # Adapters (implementations)
│       ├── adapter/
│       │   ├── http/               # HTTP handlers (driving adapter)
│       │   ├── repository/         # DB repositories (driven adapter)
│       │   ├── cache/              # Redis cache adapter
│       │   └── messaging/          # NATS messaging adapter
│       └── config/                 # Configuration
├── pkg/                            # Shared utilities
├── Dockerfile
└── go.mod
```
