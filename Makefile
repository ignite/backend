MODULE = github.com/ignite/backend
PROTO_FOLDER = api/proto
BUILD_FOLDER = bin

all: proto build

build:
	@go build -o $(BUILD_FOLDER)/ ./cmd/ignite-backend/

proto:
	@-protoc \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=module=$(MODULE) \
		--go-grpc_opt=module=$(MODULE) \
		-I $(PROTO_FOLDER) \
		$(PROTO_FOLDER)/**/*.proto

clean:
	@-rm -rf $(BUILD_FOLDER) 2> /dev/null
