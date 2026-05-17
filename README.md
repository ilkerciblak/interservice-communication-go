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

### Repository Structure

```bash
order-platform/
├── cmd/ # binaries
│   ├── gateway/        # REST → gRPC, the client-facing edge
│   ├── orderservice/   # owns orders
│   ├── inventoryservice/
│   └── notificationservice/
├── internal/
│   ├── gateway/        # RESTful API
│   ├── order/          # order domain logic
│   ├── inventory/      # inventory domain logic
│   ├── notification/   # notification domain logic
│   └── eventbus/       # messaging abstraction — empty until Stage 3
├── proto/
│   ├── order/order.proto
│   ├── inventory/inventory.proto
│   └── notification/notification.proto
├── shared/
│   ├── ... # shared utilities and types like logger,config etc
└── go.mod

```

### Project Plan and Roadmap

In summary project plan covers building interservice communication patterns with synchronous pattern first then refactoring towards asynchronous using gRPC and RabbitMQ respectively.


The system will include following system components:
- **Gateway**: exposes a RESTful API to clients. Gathers HTTP calls and handles request aggregates and internal calls.
- **Order Service**: Owns the orders domain. 
- **Inventory Service**: Handles stock opeartions.
- **Notification Service**: Handles user notification operations.

> [!NOTE]
> Further implementations as `payment, authentication & authorization` is neglected in order to focus over only `intersevice communication mechanism`. 


- **Stage 0**: Setup, protobuf, domain skeleton (~25 min): Monorepo layout, the .proto contracts, code generation. No running services yet.
- **Stage 1**: The gateway and one gRPC service (~35 min): REST gateway → gRPC → order service. The client-facing edge works end to end.
- **Stage 2**: Synchronous gRPC fan-out (~45 min): Order service calls inventory and notification over gRPC, in-line. The system works — and you'll feel exactly why this hurts.
- **Stage 3**: The messaging abstraction layer (~45 min): Define EventBus, Publisher, Subscriber interfaces. Centralized, modular, interface-driven. No behavior change yet — this is the seam.
- **Stage 4**: Migrate notification to async (~50 min): Order service publishes order.placed; notification subscribes. Drop the gRPC call. One service, fully event-driven.
- **Stage 5**: Migrate inventory to async (~40 min): Same move for inventory, plus the interesting wrinkle — inventory produces a result, so it publishes stock.reserved back. You meet event choreography.
- **Stage 6**: Production concerns (~50 min): Manual ack, dead-letter queue, graceful shutdown. The system becomes resilient instead of merely functional.

_PS: This instruction was generated via some clever AI agent **unfortunately**._


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






