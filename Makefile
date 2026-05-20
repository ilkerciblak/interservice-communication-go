proto-gen:
	@protoc --go_out=./backend/. --go-grpc_out=./backend/. backend/proto/*/*.proto

build:
	@go build -C ./backend/internal/gateway/. -o ../../cmd/gateway/ && \
		go build -C ./backend/internal/inventory/. -o ../../cmd/inventory/ && \
		go build -C ./backend/internal/notification/. -o ../../cmd/notification/ && \
		go build -C ./backend/internal/order-management/. -o ../../cmd/order-management/

run: build
	@( \
		DATA_DIR=./backend/data ./backend/cmd/gateway/gateway & \
		DATA_DIR=./backend/data ./backend/cmd/notification/notification & \
		DATA_DIR=./backend/data ./backend/cmd/inventory/inventory & \
		DATA_DIR=./backend/data ./backend/cmd/order-management/order-management & \
		trap 'kill 0' EXIT; \
		wait \
	)
