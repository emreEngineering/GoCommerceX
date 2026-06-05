# GoCommerceX

GoCommerceX is a learning-focused commerce backend platform built with Go, Hexagonal Architecture, Microservices, and gRPC.

## Purpose

The purpose of this project is to learn and practice modern backend concepts through a realistic commerce platform.

Main learning goals:

- Go backend development
- Hexagonal Architecture
- Microservice Architecture
- gRPC communication
- Protocol Buffers
- PostgreSQL
- Redis
- RabbitMQ
- JWT authentication
- Docker Compose
- Testing, logging, and health checks

## Services

- API Gateway
- Auth Service
- User Service
- Product Service
- Inventory Service
- Cart Service
- Order Service
- Payment Service
- Notification Service

## Current Phase

Phase 01: Project Foundation

Completed steps:

- 01.00 Baseline
- 01.01 Monorepo Structure
- 01.02 Docker Compose Foundation
- 01.03 Proto Foundation

Current step:

- 01.04 Project Documentation Foundation

## Project Structure

```text
GoCommerceX/
├── api-gateway/
├── auth-service/
├── user-service/
├── product-service/
├── inventory-service/
├── cart-service/
├── order-service/
├── payment-service/
├── notification-service/
├── proto/
├── deployments/
├── docs/
├── scripts/
├── .env.example
├── go.mod
└── task.md
```

## Infrastructure

The local infrastructure is defined in:

```text
deployments/docker-compose.yml
```

It includes:

- PostgreSQL
- Redis
- RabbitMQ with Management UI

To validate the Docker Compose file:

```bash
docker compose --env-file .env.example -f deployments/docker-compose.yml config
```

## Git Progress Tags

- `phase-01.00-baseline`
- `phase-01.01-monorepo-structure`
- `phase-01.02-docker-compose-foundation`
- `phase-01.03-proto-foundation`

## Development Rule

Each service owns its own data.

A service must not connect directly to another service's database. Communication between services should happen through gRPC or events.
