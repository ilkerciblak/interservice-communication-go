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

### Initial Project Plan

Initial project plan visualizes project development plan. In order to observe the differences and trade-offs between synchronous and asynchronous communication patterns interservice communications will be developed using synchronous communication pattern (with gRPC) at first. Then, interservice communications will be migrated to asynchronous communication (with RabbitMQ) using `pub/sub` pattern.

The client always interacts with the system through a RESTful API exposed by the Gateway, currently covering GET /inventory for browsing stock and POST /orders for placing new orders. To support these endpoints, the Order and Inventory services continue to communicate with the Gateway via gRPC, while all cross-service workflows are handled asynchronously through events.

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

Synchronous and asynchronous are fundamental communication patterns that dictate how systems or people exchange information. Besides other trade-offs like development complexity and architectural flexibilities, core difference between these patterns is the processing with blocking vs. non-blocking behavior. While synchronous communication blocking the next process until the sender gathers its response, in asynchronous operation pattern sender just publishes an event then proceeds without waiting and processes the reply later.


### Final Project Architecture with Event Driven Communication and gRPC 


#### Application Flows
```ascii
┌─────────────────────────────────────────────────────────────────────────────┐
│ FLOW                                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│ 1. Client  ──HTTP──►  Gateway  ──gRPC──►  Order Service                     │
│ 2. Order Service creates order (status: pending)                            │
│      ├─ returns pending order via gRPC ──► Gateway ──HTTP──► Client         │
│      └─ publishes  OrderPlacedEvent                                         │
│ 3. Inventory Service consumes OrderPlacedEvent → reserveStock()             │
│      ├─ success → publishes StockReservedEvent                              │
│      └─ failure → publishes StockNotReservedEvent → Order: rejected         │
│ 4. Notification Service consumes StockReservedEvent → notifies user         │
│      ├─ user confirms → publishes OrderConfirmedEvent → Order: confirmed    │
│      └─ user cancels  → publishes OrderCancelledEvent → Order: cancelled    │
└─────────────────────────────────────────────────────────────────────────────┘
```



```ascii
                                  ┌──────────────┐
                                  │    CLIENT    │
                                  └──────┬───────┘
                                         │
                                       HTTP
                                         │
                                         ▼
                                  ┌──────────────┐
                                  │   GATEWAY    │
                                  └──────┬───────┘
                                         │
                            ┌────────────┴────────────┐
                          gRPC                       gRPC
                            │                         │
                            ▼                         ▼
                   ┌─────────────────┐       ┌──────────────────┐
                   │  ORDER SERVICE  │       │ INVENTORY SVC    │
                   │                 │       │                  │
                   │ • create(pending)│      │ • get inventory  │
                   │ • reject()      │       │ • reserveStock() │
                   │ • confirm()     │       │                  │
                   │ • cancel()      │       │                  │
                   └────┬───────▲────┘       └────┬─────────▲───┘
                        │       │                 │         │
                        │       │                 │         │
                publishes    subscribes      subscribes  publishes
                        │       │                 │         │
                        │       │                 │         │
   ┌────────────────────┼───────┼─────────────────┼─────────┼────────────────────┐
   │                    │       │   EVENT BUS     │         │                    │
   │                    ▼       │                 ▼         │                    │
   │            OrderPlacedEvent┼─────────────────┘         │                    │
   │                            │                           │                    │
   │            StockNotReserved┘◄──────── publishes ───────┤                    │
   │                                                        │                    │
   │            StockReservedEvent ◄──────── publishes ─────┘                    │
   │                    │                                                        │
   │                    │  subscribes                                            │
   │                    ▼                                                        │
   │            ┌──────────────────────┐                                         │
   │            │ NOTIFICATION SERVICE │                                         │
   │            │                      │                                         │
   │            │ • notify user        │                                         │
   │            │ • await response     │                                         │
   │            └──────────┬───────────┘                                         │
   │                       │                                                     │
   │                   publishes                                                 │
   │                       │                                                     │
   │         ┌─────────────┴──────────────┐                                      │
   │         ▼                            ▼                                      │
   │  OrderConfirmedEvent         OrderCancelledEvent                            │
   │         │                            │                                      │
   │         └──────────┬─────────────────┘                                      │
   │                    │                                                        │
   │              subscribed by Order Service ──────────────────────────────────►│
   │                                                                             │
   └─────────────────────────────────────────────────────────────────────────────┘
```

### Repository Structure


```bash
.
├── backend
│   ├── cmd # project binaries 
│   ├── data # project data source
│   │   ├── inventory_2.json
│   │   └── order.json
│   ├── internal 
│   │   ├── event-bus
│   │   │   └── events.go # typed events and utilities
│   │   ├── gateway # RESTful API gateway
│   │   │   └── main.go 
│   │   ├── inventory # inventory domain logic
│   │   │   ├── main.go
│   │   │   ├── repository.go
│   │   │   ├── server.go
│   │   │   ├── service.go
│   │   │   └── types.go
│   │   ├── notification # notification domain logic
│   │   │   ├── main.go
│   │   │   └── server_grpc.go
│   │   └── order-management # main application logic, order-management
│   │       ├── main.go
│   │       ├── repository.go
│   │       ├── server.go
│   │       ├── service.go
│   │       └── types.go
│   ├── proto # protobuf contracts 
│   │   ├── inventory
│   │   │   └── inventory.proto
│   │   ├── notification
│   │   │   └── notification.proto
│   │   └── order
│   │       └── order.proto
│   └── shared # shared application logic, abstractions and utilities
│       ├── messaging
│       │   ├── in_memory.go
│       │   ├── rabbit.go
│       │   └── types.go
│       └── proto
│           ├── inventory
│           ├── notification
│           └── order
│   ├── go.mod
│   ├── go.sum
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md

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


- [x] **Stage 0**: Setup, protobuf, domain skeleton (~25 min): Monorepo layout, the .proto contracts, code generation. No running services yet.
- [x] **Stage 1**: The gateway and one gRPC service (~35 min): REST gateway → gRPC → order service. The client-facing edge works end to end.
- [x] **Stage 2**: Synchronous gRPC fan-out (~45 min): Order service calls inventory and notification over gRPC, in-line. The system works — and you'll feel exactly why this hurts.
- [x] **Stage 3**: The messaging abstraction layer (~45 min): Define EventBus, Publisher, Subscriber interfaces. Centralized, modular, interface-driven. No behavior change yet — this is the seam.
- [x] **Stage 4**: Migrate notification to async (~50 min): Order service publishes order.placed; notification subscribes. Drop the gRPC call. One service, fully event-driven.
- [x] **Stage 5**: Migrate inventory to async (~40 min): Same move for inventory, plus the interesting wrinkle — inventory produces a result, so it publishes stock.reserved back. You meet event choreography.
- [] **Stage 6**: Production concerns (~50 min): Manual ack, dead-letter queue, graceful shutdown. The system becomes resilient instead of merely functional.

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






