# Monitoring Platform Architecture

## Overview

Monitoring Platform is a personal information monitoring SaaS.

The system collects information from external sources,
analyzes it and sends notifications to users.

## Architecture Style

Microservice architecture.

Communication:

- REST API
- gRPC (future)
- Event-driven messaging (future)

## Initial Services

### API Gateway

Responsibilities:

- HTTP API entry point
- Authentication middleware
- Request routing


### Auth Service

Responsibilities:

- User registration
- Login
- JWT tokens


### User Service

Responsibilities:

- User profile
- Preferences


### Watcher Service

Responsibilities:

- User monitoring rules
- Interests
- Sources


## Infrastructure

Initial:

- PostgreSQL
- Redis
- Docker


Future:

- NATS/Kafka
- Elasticsearch
- ClickHouse
- Kubernetes
