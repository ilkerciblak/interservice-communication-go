proto-gen:
	@protoc --go_out=./backend/. --go_grpc_out=./backend/. backend/proto/*/*.proto
