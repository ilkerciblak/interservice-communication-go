proto-gen:
	@protoc --go_out=./backend/. --go-grpc_out=./backend/. backend/proto/*/*.proto
