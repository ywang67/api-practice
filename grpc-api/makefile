.PHONY: build_proto clean dev-server dev-client

PROTO_SRC_DIR = proto
PROTO_OUT_DIR = gen
PROTO_FILES = $(shell find $(PROTO_SRC_DIR) -name "*.proto")

build_proto:
	@echo "Generating Go code..."
	mkdir -p $(PROTO_OUT_DIR)
	protoc --proto_path=$(PROTO_SRC_DIR) \
		--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Done."

clean:
	rm -rf $(PROTO_OUT_DIR)

dev-server:
	@echo "Starting gRPC server with auto-reload..."
	reflex -r '^(server|methods|proto|gen)/.*\.go$$' --start-service -- go run server/*.go

dev-client:
	@echo "Starting gRPC client with auto-reload..."
	reflex -r '^client/.*\.go$$' --start-service -- go run client/*.go
