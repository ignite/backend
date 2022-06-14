MODULE = github.com/ignite-hq/backend
PROTO_FILES = api/proto/*
BUILD_FOLDER = bin

all: proto build

build:
	@go build -o $(BUILD_FOLDER)/ ./cmd/ignite-backend/

proto:
	@-protoc --go_out=. --go-grpc_out=. --go_opt=module=$(MODULE) --go-grpc_opt=module=$(MODULE) $(PROTO_FILES)

clean:
	@-rm -rf $(BUILD_FOLDER) 2> /dev/null
