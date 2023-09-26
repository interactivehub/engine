.PHONY: protoc
protoc:
	protoc --go_out=. --go-grpc_out=. postman/users/users_service.proto
# Add other services	