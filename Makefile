BUILD_FOLDER = bin

build:
	@go build -o $(BUILD_FOLDER)/ ./cmd/blockchain-backend/

clean:
	@-rm -rf $(BUILD_FOLDER) 2> /dev/null
