# Interservice Communication in Microservice Applications

Basic microservice _order-processing platform_ implementation that involves three decoupled services behing a RESTful API gateway. This implementation mainly scopes exercising different interservice communication patterns that provides synchronous and asynchronous via gRPC and message brookers respectively. Implementation plan consists migrating inter-service communication from synchronous gRPC to asynchronous via multiple steps.

## Tech Stack

- **Docker**: application and development environment containerization.
- **Go 1.22+**: server development.
- **gRPC-go**: `rpc` api development.
- **protocol-buffers-cli**: Code generation from `.proto` contracts.
- **RabbitMQ**: message broker. 

_Running RabbitMQ in local Docker means no signup, no credit card, no free-tier clock. Hence it is a good choice for educational purposes._



## Project Overview

```
REST               HTTP                 gRPC
   [ Client ] ───────────▶ [ Gateway ] ──────▶ [ Order Service ]
                                                       │
                                                       │  (Phase 3: sync gRPC call)
                                                       │  (Phase 5: async event)
                                                       ▼
                                              [ Notification Service ]
                                              [ Inventory Service ]
```

## Requisities and Installation

### Prerequisities

- **Docker**: In order to build and run application container. 
- **Go 1.22+** _optional_: Despite app can be run over binary, local go-cli can be used. 

1. Clone the project
```bash
go clone https://github.com/ilkerciblak/interservice-communication-go.git ~/project-dir 
```


2. Build project container 
```bash
cd ~/project-dir && \
docker compose up
```

3. Run the project
```bash
# in project directory
docker exec api make run
```






